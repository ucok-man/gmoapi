-- +goose Up
-- +goose StatementBegin
INSERT INTO movies (title, year, runtime, genres) VALUES
('The Shawshank Redemption', 1994, 142, ARRAY['Drama']),
('The Godfather', 1972, 175, ARRAY['Crime', 'Drama']),
('The Dark Knight', 2008, 152, ARRAY['Action', 'Crime', 'Drama']),
('Pulp Fiction', 1994, 154, ARRAY['Crime', 'Drama']),
('The Lord of the Rings: The Return of the King', 2003, 201, ARRAY['Adventure', 'Drama', 'Fantasy']),
('Forrest Gump', 1994, 142, ARRAY['Drama', 'Romance']),
('Inception', 2010, 148, ARRAY['Action', 'Adventure', 'Sci-Fi']),
('The Matrix', 1999, 136, ARRAY['Action', 'Sci-Fi']),
('Schindler''s List', 1993, 195, ARRAY['Biography', 'Drama', 'History']),
('Parasite', 2019, 132, ARRAY['Comedy', 'Drama', 'Thriller']),
('Goodfellas', 1990, 146, ARRAY['Biography', 'Crime', 'Drama']),
('Spirited Away', 2001, 125, ARRAY['Animation', 'Adventure', 'Family']),
('Fight Club', 1999, 139, ARRAY['Drama']),
('The Green Mile', 1999, 189, ARRAY['Crime', 'Drama', 'Fantasy']),
('Gladiator', 2000, 155, ARRAY['Action', 'Adventure', 'Drama']),
('Interstellar', 2014, 169, ARRAY['Adventure', 'Drama', 'Sci-Fi']),
('Alien', 1979, 117, ARRAY['Horror', 'Sci-Fi']),
('E.T. the Extra-Terrestrial', 1982, 115, ARRAY['Family', 'Sci-Fi']),
('Mad Max: Fury Road', 2015, 120, ARRAY['Action', 'Adventure', 'Sci-Fi']),
('Coco', 2017, 105, ARRAY['Animation', 'Adventure', 'Family']),

('The Silence of the Lambs', 1991, 118, ARRAY['Crime', 'Drama', 'Thriller']),
('Se7en', 1995, 127, ARRAY['Crime', 'Drama', 'Mystery']),
('The Usual Suspects', 1995, 106, ARRAY['Crime', 'Drama', 'Mystery']),
('Saving Private Ryan', 1998, 169, ARRAY['Drama', 'War']),
('The Lion King', 1994, 88, ARRAY['Animation', 'Adventure', 'Drama']),
('Back to the Future', 1985, 116, ARRAY['Adventure', 'Comedy', 'Sci-Fi']),
('The Prestige', 2006, 130, ARRAY['Drama', 'Mystery', 'Sci-Fi']),
('Whiplash', 2014, 106, ARRAY['Drama', 'Music']),
('The Departed', 2006, 151, ARRAY['Crime', 'Drama', 'Thriller']),
('City of God', 2002, 130, ARRAY['Crime', 'Drama']),
('The Pianist', 2002, 150, ARRAY['Biography', 'Drama', 'Music']),
('Avengers: Endgame', 2019, 181, ARRAY['Action', 'Adventure', 'Drama']),
('Titanic', 1997, 195, ARRAY['Drama', 'Romance']),
('Shutter Island', 2010, 138, ARRAY['Mystery', 'Thriller']),
('The Social Network', 2010, 120, ARRAY['Biography', 'Drama']),
('Joker', 2019, 122, ARRAY['Crime', 'Drama', 'Thriller']),
('The Grand Budapest Hotel', 2014, 99, ARRAY['Adventure', 'Comedy', 'Drama']),
('WALL·E', 2008, 98, ARRAY['Animation', 'Adventure', 'Family']),
('Inside Out', 2015, 95, ARRAY['Animation', 'Adventure', 'Comedy']),
('A Beautiful Mind', 2001, 135, ARRAY['Biography', 'Drama']),

('Apocalypse Now', 1979, 147, ARRAY['Drama', 'War']),
('Casablanca', 1942, 102, ARRAY['Drama', 'Romance', 'War']),
('Citizen Kane', 1941, 119, ARRAY['Drama', 'Mystery']),
('Lawrence of Arabia', 1962, 222, ARRAY['Adventure', 'Biography', 'Drama']),
('Psycho', 1960, 109, ARRAY['Horror', 'Mystery', 'Thriller']),
('2001: A Space Odyssey', 1968, 149, ARRAY['Adventure', 'Sci-Fi']),
('The Good, the Bad and the Ugly', 1966, 178, ARRAY['Western']),
('Dr. Strangelove', 1964, 95, ARRAY['Comedy', 'War']),
('Singin'' in the Rain', 1952, 103, ARRAY['Comedy', 'Musical', 'Romance']),
('It''s a Wonderful Life', 1946, 130, ARRAY['Drama', 'Family', 'Fantasy']),
('Metropolis', 1927, 153, ARRAY['Drama', 'Sci-Fi']),
('Rear Window', 1954, 112, ARRAY['Mystery', 'Thriller']),
('North by Northwest', 1959, 136, ARRAY['Adventure', 'Mystery', 'Thriller']),
('Some Like It Hot', 1959, 121, ARRAY['Comedy', 'Romance']),
('The Wizard of Oz', 1939, 102, ARRAY['Adventure', 'Family', 'Fantasy']),
('Sunset Boulevard', 1950, 110, ARRAY['Drama', 'Film-Noir']),
('12 Angry Men', 1957, 96, ARRAY['Drama']),
('Gone with the Wind', 1939, 221, ARRAY['Drama', 'Romance', 'War']),
('Chinatown', 1974, 131, ARRAY['Drama', 'Mystery', 'Thriller']),
('Rashomon', 1950, 88, ARRAY['Crime', 'Drama', 'Mystery']),

('Dune', 2021, 155, ARRAY['Adventure', 'Drama', 'Sci-Fi']),
('Everything Everywhere All at Once', 2022, 139, ARRAY['Action', 'Adventure', 'Comedy']),
('The Batman', 2022, 176, ARRAY['Action', 'Crime', 'Drama']),
('Oppenheimer', 2023, 180, ARRAY['Biography', 'Drama', 'History']),
('Barbie', 2023, 114, ARRAY['Adventure', 'Comedy', 'Fantasy']),
('Spider-Man: Into the Spider-Verse', 2018, 117, ARRAY['Animation', 'Action', 'Adventure']),
('Spider-Man: Across the Spider-Verse', 2023, 140, ARRAY['Animation', 'Action', 'Adventure']),
('The Irishman', 2019, 209, ARRAY['Biography', 'Crime', 'Drama']),
('Soul', 2020, 100, ARRAY['Animation', 'Adventure', 'Comedy']),
('Inside Llewyn Davis', 2013, 104, ARRAY['Drama', 'Music']),
('La La Land', 2016, 128, ARRAY['Comedy', 'Drama', 'Music']),
('The Revenant', 2015, 156, ARRAY['Action', 'Adventure', 'Drama']),
('Birdman', 2014, 119, ARRAY['Comedy', 'Drama']),
('Her', 2013, 126, ARRAY['Drama', 'Romance', 'Sci-Fi']),
('Moonlight', 2016, 111, ARRAY['Drama']),
('The Shape of Water', 2017, 123, ARRAY['Adventure', 'Drama', 'Fantasy']),
('Jojo Rabbit', 2019, 108, ARRAY['Comedy', 'Drama', 'War']),
('Black Panther', 2018, 134, ARRAY['Action', 'Adventure', 'Sci-Fi']),
('The Whale', 2022, 117, ARRAY['Drama']),
('Tenet', 2020, 150, ARRAY['Action', 'Sci-Fi', 'Thriller']),

('Toy Story', 1995, 81, ARRAY['Animation', 'Adventure', 'Comedy']),
('Toy Story 3', 2010, 103, ARRAY['Animation', 'Adventure', 'Comedy']),
('Finding Nemo', 2003, 100, ARRAY['Animation', 'Adventure', 'Comedy']),
('Finding Dory', 2016, 97, ARRAY['Animation', 'Adventure', 'Comedy']),
('Up', 2009, 96, ARRAY['Animation', 'Adventure', 'Comedy']),
('Ratatouille', 2007, 111, ARRAY['Animation', 'Adventure', 'Comedy']),
('Monsters, Inc.', 2001, 92, ARRAY['Animation', 'Adventure', 'Comedy']),
('Monsters University', 2013, 104, ARRAY['Animation', 'Adventure', 'Comedy']),
('Shrek', 2001, 90, ARRAY['Animation', 'Adventure', 'Comedy']),
('Shrek 2', 2004, 93, ARRAY['Animation', 'Adventure', 'Comedy']),
('Kung Fu Panda', 2008, 92, ARRAY['Animation', 'Action', 'Adventure']),
('Kung Fu Panda 2', 2011, 91, ARRAY['Animation', 'Action', 'Adventure']),
('Frozen', 2013, 102, ARRAY['Animation', 'Adventure', 'Comedy']),
('Frozen II', 2019, 103, ARRAY['Animation', 'Adventure', 'Comedy']),
('Zootopia', 2016, 108, ARRAY['Animation', 'Adventure', 'Comedy']),
('Moana', 2016, 107, ARRAY['Animation', 'Adventure', 'Comedy']),
('Encanto', 2021, 102, ARRAY['Animation', 'Comedy', 'Family']),
('Despicable Me', 2010, 95, ARRAY['Animation', 'Adventure', 'Comedy']),
('Despicable Me 2', 2013, 98, ARRAY['Animation', 'Adventure', 'Comedy']),
('How to Train Your Dragon', 2010, 98, ARRAY['Animation', 'Action', 'Adventure']);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movies WHERE 
(title = 'The Shawshank Redemption' AND year = 1994) OR
(title = 'The Godfather' AND year = 1972) OR
(title = 'The Dark Knight' AND year = 2008) OR
(title = 'Pulp Fiction' AND year = 1994) OR
(title = 'The Lord of the Rings: The Return of the King' AND year = 2003) OR
(title = 'Forrest Gump' AND year = 1994) OR
(title = 'Inception' AND year = 2010) OR
(title = 'The Matrix' AND year = 1999) OR
(title = 'Schindler''s List' AND year = 1993) OR
(title = 'Parasite' AND year = 2019) OR
(title = 'Goodfellas' AND year = 1990) OR
(title = 'Spirited Away' AND year = 2001) OR
(title = 'Fight Club' AND year = 1999) OR
(title = 'The Green Mile' AND year = 1999) OR
(title = 'Gladiator' AND year = 2000) OR
(title = 'Interstellar' AND year = 2014) OR
(title = 'Alien' AND year = 1979) OR
(title = 'E.T. the Extra-Terrestrial' AND year = 1982) OR
(title = 'Mad Max: Fury Road' AND year = 2015) OR
(title = 'Coco' AND year = 2017) OR

(title = 'The Silence of the Lambs' AND year = 1991) OR
(title = 'Se7en' AND year = 1995) OR
(title = 'The Usual Suspects' AND year = 1995) OR
(title = 'Saving Private Ryan' AND year = 1998) OR
(title = 'The Lion King' AND year = 1994) OR
(title = 'Back to the Future' AND year = 1985) OR
(title = 'The Prestige' AND year = 2006) OR
(title = 'Whiplash' AND year = 2014) OR
(title = 'The Departed' AND year = 2006) OR
(title = 'City of God' AND year = 2002) OR
(title = 'The Pianist' AND year = 2002) OR
(title = 'Avengers: Endgame' AND year = 2019) OR
(title = 'Titanic' AND year = 1997) OR
(title = 'Shutter Island' AND year = 2010) OR
(title = 'The Social Network' AND year = 2010) OR
(title = 'Joker' AND year = 2019) OR
(title = 'The Grand Budapest Hotel' AND year = 2014) OR
(title = 'WALL·E' AND year = 2008) OR
(title = 'Inside Out' AND year = 2015) OR
(title = 'A Beautiful Mind' AND year = 2001) OR

(title = 'Apocalypse Now' AND year = 1979) OR
(title = 'Casablanca' AND year = 1942) OR
(title = 'Citizen Kane' AND year = 1941) OR
(title = 'Lawrence of Arabia' AND year = 1962) OR
(title = 'Psycho' AND year = 1960) OR
(title = '2001: A Space Odyssey' AND year = 1968) OR
(title = 'The Good, the Bad and the Ugly' AND year = 1966) OR
(title = 'Dr. Strangelove' AND year = 1964) OR
(title = 'Singin'' in the Rain' AND year = 1952) OR
(title = 'It''s a Wonderful Life' AND year = 1946) OR
(title = 'Metropolis' AND year = 1927) OR
(title = 'Rear Window' AND year = 1954) OR
(title = 'North by Northwest' AND year = 1959) OR
(title = 'Some Like It Hot' AND year = 1959) OR
(title = 'The Wizard of Oz' AND year = 1939) OR
(title = 'Sunset Boulevard' AND year = 1950) OR
(title = '12 Angry Men' AND year = 1957) OR
(title = 'Gone with the Wind' AND year = 1939) OR
(title = 'Chinatown' AND year = 1974) OR
(title = 'Rashomon' AND year = 1950) OR

(title = 'Dune' AND year = 2021) OR
(title = 'Everything Everywhere All at Once' AND year = 2022) OR
(title = 'The Batman' AND year = 2022) OR
(title = 'Oppenheimer' AND year = 2023) OR
(title = 'Barbie' AND year = 2023) OR
(title = 'Spider-Man: Into the Spider-Verse' AND year = 2018) OR
(title = 'Spider-Man: Across the Spider-Verse' AND year = 2023) OR
(title = 'The Irishman' AND year = 2019) OR
(title = 'Soul' AND year = 2020) OR
(title = 'Inside Llewyn Davis' AND year = 2013) OR
(title = 'La La Land' AND year = 2016) OR
(title = 'The Revenant' AND year = 2015) OR
(title = 'Birdman' AND year = 2014) OR
(title = 'Her' AND year = 2013) OR
(title = 'Moonlight' AND year = 2016) OR
(title = 'The Shape of Water' AND year = 2017) OR
(title = 'Jojo Rabbit' AND year = 2019) OR
(title = 'Black Panther' AND year = 2018) OR
(title = 'The Whale' AND year = 2022) OR
(title = 'Tenet' AND year = 2020) OR

(title = 'Toy Story' AND year = 1995) OR
(title = 'Toy Story 3' AND year = 2010) OR
(title = 'Finding Nemo' AND year = 2003) OR
(title = 'Finding Dory' AND year = 2016) OR
(title = 'Up' AND year = 2009) OR
(title = 'Ratatouille' AND year = 2007) OR
(title = 'Monsters, Inc.' AND year = 2001) OR
(title = 'Monsters University' AND year = 2013) OR
(title = 'Shrek' AND year = 2001) OR
(title = 'Shrek 2' AND year = 2004) OR
(title = 'Kung Fu Panda' AND year = 2008) OR
(title = 'Kung Fu Panda 2' AND year = 2011) OR
(title = 'Frozen' AND year = 2013) OR
(title = 'Frozen II' AND year = 2019) OR
(title = 'Zootopia' AND year = 2016) OR
(title = 'Moana' AND year = 2016) OR
(title = 'Encanto' AND year = 2021) OR
(title = 'Despicable Me' AND year = 2010) OR
(title = 'Despicable Me 2' AND year = 2013) OR
(title = 'How to Train Your Dragon' AND year = 2010);
-- +goose StatementEnd
