# Book Store RESTful API

RESTful API untuk sistem toko buku dengan fitur authentication, manajemen produk, keranjang belanja, dan pemesanan.

## Struktur Folder

```
.
├── config/          # Konfigurasi database
├── controller/      # HTTP handlers
├── entity/          # Domain entities
├── middleware/      # HTTP middleware (auth, cors)
├── model/           # Request/Response DTOs
├── repository/      # Database layer
├── router/          # Route definitions
├── service/         # Business logic
└── main.go         # Application entry point
```

## Tech Stack

- Go 1.21+
- PostgreSQL
- Native `net/http` (no framework)
- Clean Architecture

## Setup

### 1. Environment Variables

Buat file `.env`:

```env
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=5432
DB_NAME=bookstore
DB_SSLMODE=disable
PORT=8080
BASE_URL=http://localhost:8080
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run Application

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`

## API Documentation

### Authentication

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john",
  "password": "password123",
  "email": "john@example.com"
}
```

Response:
```json
{
  "status": "success",
  "message": "User registered successfully"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

Response:
```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "abc123...",
    "username": "admin",
    "role": "admin"
  }
}
```

#### Logout
```http
POST /api/auth/logout
Authorization: Bearer {token}
```

### Books

#### Get All Books
```http
GET /api/books
```

Response:
```json
{
  "status": "success",
  "message": "Books retrieved successfully",
  "data": [
    {
      "id": 1,
      "nama_barang": "Go Programming",
      "stok": 10,
      "terjual": 5,
      "harga": 150000,
      "keterangan": "Book about Go",
      "gambar_buku": "http://localhost:8080/uploads/books/1234567890_abc123.jpg",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### Get Book by ID
```http
GET /api/books/detail?id=1
```

#### Create Book (Admin Only)

**Option 1: With Image Upload (multipart/form-data)**
```http
POST /api/books
Authorization: Bearer {admin_token}
Content-Type: multipart/form-data

Form Data:
- nama_barang: "Go Programming"
- stok: 10
- harga: 150000
- keterangan: "Book about Go programming"
- gambar_buku: [file upload]
```

**Option 2: Without Image (multipart/form-data)**
```http
POST /api/books
Authorization: Bearer {admin_token}
Content-Type: multipart/form-data

Form Data:
- nama_barang: "Go Programming"
- stok: 10
- harga: 150000
- keterangan: "Book about Go programming"
```

Response:
```json
{
  "status": "success",
  "message": "Book created successfully",
  "data": {
    "id": 1,
    "nama_barang": "Go Programming",
    "stok": 10,
    "terjual": 0,
    "harga": 150000,
    "keterangan": "Book about Go",
    "gambar_buku": "http://localhost:8080/uploads/books/1234567890_abc123.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Update Book (Admin Only)
```http
PUT /api/books/detail?id=1
Authorization: Bearer {admin_token}
Content-Type: multipart/form-data

Form Data:
- nama_barang: "Go Programming Advanced"
- stok: 15
- terjual: 5
- harga: 175000
- keterangan: "Advanced Go book"
- gambar_buku: [file upload] (optional, only if changing image)
```

Response:
```json
{
  "status": "success",
  "message": "Book updated successfully",
  "data": {
    "id": 1,
    "nama_barang": "Go Programming Advanced",
    "stok": 15,
    "terjual": 5,
    "harga": 175000,
    "keterangan": "Advanced Go book",
    "gambar_buku": "http://localhost:8080/uploads/books/1234567890_xyz789.jpg",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-02T00:00:00Z"
  }
}
```

#### Delete Book (Admin Only)
```http
DELETE /api/books/detail?id=1
Authorization: Bearer {admin_token}
```

### Cart

#### Get Cart
```http
GET /api/cart
Authorization: Bearer {token}
```

Response:
```json
{
  "status": "success",
  "message": "Cart retrieved successfully",
  "data": {
    "items": [
      {
        "id": 1,
        "book_id": 1,
        "nama_barang": "Go Programming",
        "jumlah": 2,
        "harga": 150000,
        "stok": 10,
        "harga_satuan": 150000,
        "subtotal": 300000,
        "gambar_buku": "http://localhost:8080/uploads/books/1234567890_abc123.jpg"
      }
    ],
    "total": 300000
  }
}
```

#### Add to Cart
```http
POST /api/cart
Authorization: Bearer {token}
Content-Type: application/json

{
  "book_id": 1,
  "jumlah": 2
}
```

#### Update Cart Item
```http
PUT /api/cart/item?id=1
Authorization: Bearer {token}
Content-Type: application/json

{
  "jumlah": 3
}
```

#### Remove from Cart
```http
DELETE /api/cart/item?id=1
Authorization: Bearer {token}
```

#### Clear Cart
```http
DELETE /api/cart
Authorization: Bearer {token}
```

### Orders

#### Create Order (Checkout)
```http
POST /api/orders
Authorization: Bearer {token}
```

Response:
```json
{
  "status": "success",
  "message": "Order created successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "total_harga": 300000,
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Get User Orders
```http
GET /api/orders
Authorization: Bearer {token}
```

#### Get Order Detail
```http
GET /api/orders/detail?id=1
Authorization: Bearer {token}
```

Response:
```json
{
  "status": "success",
  "message": "Order detail retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 2,
    "username": "john",
    "total_harga": 300000,
    "status": "pending",
    "created_at": "2024-01-01T00:00:00Z",
    "items": [
      {
        "id": 1,
        "order_id": 1,
        "book_id": 1,
        "jumlah": 2,
        "harga": 150000,
        "created_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

### Health Check
```http
GET /api/health
```

## Default Users

Aplikasi sudah menyediakan 2 user default:

1. **Admin**
   - Username: `admin`
   - Password: `admin123`
   - Role: `admin`

2. **User**
   - Username: `user`
   - Password: `user123`
   - Role: `user`

## Authorization

### User Roles

- **Admin**: Dapat melakukan CRUD pada produk
- **User**: Hanya dapat membeli (cart & order)

### Protected Endpoints

Endpoint yang memerlukan authentication menggunakan Bearer token di header:

```
Authorization: Bearer {your_token_here}
```

Token didapatkan dari endpoint login.

## Response Format

Semua response menggunakan format standar:

### Success Response
```json
{
  "status": "success",
  "message": "Operation successful",
  "data": {}
}
```

### Error Response
```json
{
  "status": "error",
  "message": "Error description"
}
```

## Testing dengan cURL

### Login sebagai Admin
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### Create Book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "nama_barang=Go Programming" \
  -F "stok=10" \
  -F "harga=150000" \
  -F "keterangan=Learn Go" \
  -F "gambar_buku=@/path/to/image.jpg"
```

### Create Book Without Image
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "nama_barang=Go Programming" \
  -F "stok=10" \
  -F "harga=150000" \
  -F "keterangan=Learn Go"
```

### Update Book with New Image
```bash
curl -X PUT "http://localhost:8080/api/books/detail?id=1" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "nama_barang=Go Programming Advanced" \
  -F "stok=15" \
  -F "terjual=5" \
  -F "harga=175000" \
  -F "keterangan=Advanced Go" \
  -F "gambar_buku=@/path/to/new-image.jpg"
```

### Upload Image Separately
```bash
curl -X POST http://localhost:8080/api/upload/image \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -F "image=@/path/to/image.jpg"
```

### Get All Books
```bash
curl http://localhost:8080/api/books
```

## Testing dengan Postman

### 1. Login dan Dapatkan Token

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/auth/login`
- Body (JSON):
```json
{
  "username": "admin",
  "password": "admin123"
}
```

Copy token dari response untuk digunakan di request selanjutnya.

### 2. Create Book dengan Upload Image

**Request:**
- Method: `POST`
- URL: `http://localhost:8080/api/books`
- Headers:
  - `Authorization: Bearer YOUR_TOKEN`
- Body (form-data):
  - `nama_barang`: Go Programming
  - `stok`: 10
  - `harga`: 150000
  - `keterangan`: Learn Go from basics
  - `gambar_buku`: [Select File]

### 3. Update Book dengan Ganti Image

**Request:**
- Method: `PUT`
- URL: `http://localhost:8080/api/books/detail?id=1`
- Headers:
  - `Authorization: Bearer YOUR_TOKEN`
- Body (form-data):
  - `nama_barang`: Go Programming Advanced
  - `stok`: 15
  - `terjual`: 5
  - `harga`: 175000
  - `keterangan`: Advanced topics
  - `gambar_buku`: [Select New File] (optional)

### 4. Get Book dengan Image URL

**Request:**
- Method: `GET`
- URL: `http://localhost:8080/api/books/detail?id=1`

**Response akan include full image URL:**
```json
{
  "status": "success",
  "message": "Book retrieved successfully",
  "data": {
    "id": 1,
    "gambar_buku": "http://localhost:8080/uploads/books/1234567890_abc123.jpg"
  }
}
```

### 5. Access Image di Browser

Buka URL image langsung di browser:
```
http://localhost:8080/uploads/books/1234567890_abc123.jpg
```

## Image Storage

**Folder Structure:**
```
uploads/
└── books/
    ├── 1234567890_abc123.jpg
    ├── 1234567891_def456.png
    └── 1234567892_ghi789.jpg
```

**Features:**
- Automatic unique filename generation (timestamp + random string)
- File validation (type & size)
- Automatic old image deletion when updating
- Image cleanup when deleting books
- Public access to view images
- Cache headers for better performance

## Best Practices yang Diterapkan

1. **Clean Architecture**: Pemisahan concerns dengan layer yang jelas
2. **Dependency Injection**: Loose coupling antar komponen
3. **Repository Pattern**: Abstraksi database operations
4. **Middleware Pattern**: Reusable HTTP middleware
5. **Error Handling**: Consistent error responses
6. **Validation**: Input validation di service layer
7. **Transaction Management**: Database transactions untuk operasi kompleks
8. **RESTful Design**: HTTP methods dan status codes yang sesuai

## License

MIT