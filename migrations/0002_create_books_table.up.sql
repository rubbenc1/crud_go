CREATE TABLE IF NOT EXISTS books (
    BookId SERIAL PRIMARY KEY,
    Name VARCHAR(255) NOT NULL,
    PublishedYear INT NOT NULL,
    ShelfId INT REFERENCES shelves(ShelfId)
);