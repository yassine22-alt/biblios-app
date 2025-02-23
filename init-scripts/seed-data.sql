-- Seed Authors data
INSERT INTO authors (id, first_name, last_name, bio) VALUES
(1, 'Mark', 'Twain', 'Mark Twain, born Samuel Langhorne Clemens, was an American writer and humorist best known for ''The Adventures of Tom Sawyer'' and ''Adventures of Huckleberry Finn''.'),
(2, 'Mary', 'Shelley', 'Mary Shelley was an English author best known for her Gothic novel ''Frankenstein'', which is often considered one of the first science fiction works.'),
(3, 'brs', 'blali', 'lorem epsium');

-- Seed Books data
INSERT INTO books (id, title, author_id, published_at, price, stock) VALUES
(1, 'Pride and Prejudice', 1, '2025-01-12 17:54:53.369326+01:00', 12.99, 100),
(2, '1984', 3, '2025-01-12 17:56:36.8499153+01:00', 15.50, 200);

-- Seed Book Genres
INSERT INTO book_genres (book_id, genre) VALUES
(1, 'Romance'),
(1, 'Classic'),
(2, 'Dystopian'),
(2, 'Political Fiction'),
(2, 'Science Fiction');

-- Reset sequence values to handle future inserts correctly
SELECT setval('authors_id_seq', (SELECT MAX(id) FROM authors));
SELECT setval('books_id_seq', (SELECT MAX(id) FROM books));