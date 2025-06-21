package migrations

import (
	"gorm.io/gorm"
)

func RunInitMigrations(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(
			`
			drop table if exists books_authors;
			drop table if exists books;
			drop table if exists authors;

			create table books (
			id serial primary key,
			title varchar(64) not null,
			pages integer not null,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone
			);

			create table authors (
			id serial primary key,
			name varchar(64) not null,
			created_at timestamp with time zone,
			updated_at timestamp with time zone,
			deleted_at timestamp with time zone
			);

			create table books_authors (
			book_id integer not null,
			author_id integer not null,
			primary key (book_id, author_id),
			constraint fk_book foreign key (book_id) references books(id) on delete cascade,
			constraint fk_author foreign key (author_id) references authors(id) on delete cascade
			);`).Error
	})
}
