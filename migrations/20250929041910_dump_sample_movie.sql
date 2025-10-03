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
('Coco', 2017, 105, ARRAY['Animation', 'Adventure', 'Family']);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movies WHERE (title = 'The Shawshank Redemption' AND year = 1994) OR
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
(title = 'Coco' AND year = 2017);
-- +goose StatementEnd