package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic
		// as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a panic or
			// not.
			if err := recover(); err != nil {
				// If there was a panic, set a "Connection: close" header on the
				// response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after a response has been
				// sent.
				w.Header().Set("Connection", "close")

				// The value returned by recover() has the type any, so we use
				// fmt.Errorf() to normalize it into an error and call our
				// serverErrorResponse() helper.
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

/*
Pola ini hanya berlaku bila API berjalan di **satu mesin**. Kalau infrastruktur **distributed** (misalnya API berjalan di banyak server di belakang load balancer), maka dibutuhkan pendekatan lain.

Alternatif:

- Gunakan fitur rate limiting bawaan di **HAProxy** atau **Nginx**.
- Atau gunakan database cepat seperti **Redis** untuk menyimpan counter request per-client, agar semua server bisa berbagi informasi rate limit.
*/

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the clients' IP addresses and cleint struct (limiter and lastseen)
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	if app.config.Limiter.Enabled {
		// Launch a background goroutine which removes old entries from the clients map once every minute.
		go func() {
			for {
				time.Sleep(time.Minute)
				mu.Lock()
				// Loop through all clients. If they haven't been seen within the last three minutes, delete.
				for ip, client := range clients {
					if time.Since(client.lastSeen) > 3*time.Minute {
						delete(clients, ip)
					}
				}
				mu.Unlock()
			}
		}()
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.Limiter.Enabled {
			// Use the realip.FromRequest() function to get the client's IP address.
			ip := realip.FromRequest(r)

			// Lock the mutex to prevent this code from being executed concurrently.
			mu.Lock()

			// Check to see if the IP address already exists in the map. If it doesn't, then
			// initialize a new rate limiter and add the IP address and limiter to the map.
			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.config.Limiter.Rps), app.config.Limiter.Burst)}
			}

			// Update the last seen time for the client.
			clients[ip].lastSeen = time.Now()

			// Call the Allow() method on the rate limiter for the current IP address. If
			// the request isn't allowed, unlock the mutex and send a 429 Too Many Requests.
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Unlock mutex before continuing to next handler.
			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
