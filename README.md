# Book Database
## Tech stack :
- Go 
- Postgres
- mux

## Functionalities
- View Books  : /api/books `GET`
- Update Book : /api/book/{id} `PUT`
- Delete Book : /api/books/{id} `DELETE`
- Create Book : /api/books `POST`
  - Body
  ```json
    {
      Isbn : Int
      Title : String
      Author : String
     }
  ```
