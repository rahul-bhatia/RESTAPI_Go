# RESTAPI_Go
## Tech stack :
- Go 
- Postgres

## Functionalities
- View Books  : /api/books `GET`
- Update Book : /api/book/{id} `PUT`
- Delete Book : /api/books/{id} `DELETE`
- Create Book : /api/books `POST`
  - Body
    - Isbn : Int
    - Title : String
    - Author : String
