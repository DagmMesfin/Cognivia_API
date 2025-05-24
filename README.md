 # Cognivia API Documentation

## Overview
The Cognivia API is a RESTful API built with Go and Gin framework that provides endpoints for managing users, notebooks, snapnotes, and prep pilots. The API uses JWT authentication for protected routes and MongoDB as the database.

**Base URL:** `http://localhost:8080`

## Authentication
Most endpoints require JWT authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

## Data Models

### User
```json
{
  "id": "string (ObjectID)",
  "email": "string",
  "name": "string",
  "created_at": "string (ISO 8601)",
  "updated_at": "string (ISO 8601)"
}
```

### Notebook
```json
{
  "id": "string (ObjectID)",
  "user_id": "string (ObjectID)",
  "snapnotes_id": "string (ObjectID, optional)",
  "prep_pilot_id": "string (ObjectID, optional)",
  "name": "string (required)",
  "icon": "string",
  "color": "string",
  "type": "string",
  "google_drive_link": "string (optional)",
  "created_at": "string (ISO 8601)",
  "updated_at": "string (ISO 8601)"
}
```

### Snapnotes
```json
{
  "id": "string (ObjectID)",
  "title": "string",
  "summaryByChapter": [
    {
      "chapterTitle": "string",
      "summary": "string",
      "keyPoints": ["string"]
    }
  ],
  "flashcards": [
    {
      "chapterTitle": "string",
      "flashcards": [
        {
          "key term": "string",
          "definition": "string"
        }
      ]
    }
  ]
}
```

### PrepPilot
```json
{
  "id": "string (ObjectID)",
  "notebook_id": "string (ObjectID)",
  "chapters": [
    {
      "chapterTitle": "string",
      "questions": [
        {
          "question": "string",
          "options": {
            "A": "string",
            "B": "string",
            "C": "string",
            "D": "string"
          },
          "answer": "string",
          "explanation": "string"
        }
      ]
    }
  ]
}
```

## API Endpoints

### User Management

#### Register User
- **POST** `/api/v1/users/register`
- **Authentication:** None required
- **Description:** Register a new user account

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securePassword123",
    "name": "John Doe"
  }'
```

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123",
  "name": "John Doe"
}
```

**Example Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "john.doe@example.com",
    "name": "John Doe"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or validation errors
```json
{
  "error": "Key: 'User.Email' Error:Tag 'email' validation failed"
}
```
- `500 Internal Server Error`: User already exists or server error
```json
{
  "error": "user already exists"
}
```

#### Login User
- **POST** `/api/v1/users/login`
- **Authentication:** None required
- **Description:** Authenticate user and receive JWT token

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "securePassword123"
  }'
```

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "securePassword123"
}
```

**Example Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImpvaG4uZG9lQGV4YW1wbGUuY29tIiwiZXhwIjoxNzA0MDY3MjAwLCJ1c2VyX2lkIjoiNTA3ZjFmNzdiY2Y4NmNkNzk5NDM5MDExIn0.abc123def456ghi789",
  "user": {
    "id": "507f1f77bcf86cd799439011",
    "email": "john.doe@example.com",
    "name": "John Doe"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or missing required fields
```json
{
  "error": "Key: 'loginRequest.Email' Error:Tag 'required' validation failed"
}
```
- `401 Unauthorized`: Invalid credentials
```json
{
  "error": "Invalid credentials"
}
```

#### Get User by ID
- **GET** `/api/v1/users/{id}`
- **Authentication:** None required
- **Description:** Retrieve user information by user ID

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json"
```

**Path Parameters:**
- `id` (string): User ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "email": "john.doe@example.com",
  "name": "John Doe",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**
- `404 Not Found`: User not found
```json
{
  "error": "User not found"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Invalid ObjectID format"
}
```

#### Update User
- **PUT** `/api/v1/users/{id}`
- **Authentication:** None required
- **Description:** Update user information

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.updated@example.com",
    "name": "John Updated",
    "password": "newSecurePassword123"
  }'
```

**Path Parameters:**
- `id` (string): User ID (MongoDB ObjectID)

**Request Body:**
```json
{
  "email": "john.updated@example.com",
  "name": "John Updated",
  "password": "newSecurePassword123"
}
```

**Example Response (200 OK):**
```json
{
  "message": "User updated successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body or user ID
```json
{
  "error": "Invalid user ID"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Database connection error"
}
```

#### Delete User
- **DELETE** `/api/v1/users/{id}`
- **Authentication:** None required
- **Description:** Delete a user account

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/507f1f77bcf86cd799439011 \
  -H "Content-Type: application/json"
```

**Path Parameters:**
- `id` (string): User ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

**Error Responses:**
- `500 Internal Server Error`: Server error
```json
{
  "error": "Failed to delete user"
}
```

### Notebook Management
**Note:** All notebook endpoints require JWT authentication.

#### Create Notebook
- **POST** `/api/v1/notebooks/`
- **Authentication:** Required (JWT)
- **Description:** Create a new notebook for the authenticated user

**Example Request:**
```bash
curl -X POST http://localhost:8080/api/v1/notebooks/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "Biology Study Guide",
    "icon": "🧬",
    "color": "#4CAF50",
    "type": "study",
    "google_drive_link": "https://drive.google.com/file/d/1abc123def456"
  }'
```

**Request Body:**
```json
{
  "name": "Biology Study Guide",
  "icon": "🧬",
  "color": "#4CAF50",
  "type": "study",
  "google_drive_link": "https://drive.google.com/file/d/1abc123def456"
}
```

**Example Response (201 Created):**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "user_id": "507f1f77bcf86cd799439011",
  "snapnotes_id": null,
  "prep_pilot_id": null,
  "name": "Biology Study Guide",
  "icon": "🧬",
  "color": "#4CAF50",
  "type": "study",
  "google_drive_link": "https://drive.google.com/file/d/1abc123def456",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body (name is required)
```json
{
  "error": "Key: 'Notebook.Name' Error:Tag 'required' validation failed"
}
```
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "Missing or invalid Authorization header"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Failed to create notebook"
}
```

#### Get Notebook by ID
- **GET** `/api/v1/notebooks/{id}`
- **Authentication:** Required (JWT)
- **Description:** Retrieve a specific notebook by ID (user can only access their own notebooks)

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/notebooks/507f1f77bcf86cd799439012 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Path Parameters:**
- `id` (string): Notebook ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "user_id": "507f1f77bcf86cd799439011",
  "snapnotes_id": "507f1f77bcf86cd799439013",
  "prep_pilot_id": "507f1f77bcf86cd799439014",
  "name": "Biology Study Guide",
  "icon": "🧬",
  "color": "#4CAF50",
  "type": "study",
  "google_drive_link": "https://drive.google.com/file/d/1abc123def456",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `404 Not Found`: Notebook not found
```json
{
  "error": "Notebook not found"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Database connection error"
}
```

#### Get User's Notebooks
- **GET** `/api/v1/notebooks/user`
- **Authentication:** Required (JWT)
- **Description:** Retrieve all notebooks belonging to the authenticated user

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/notebooks/user \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Example Response (200 OK):**
```json
[
  {
    "id": "507f1f77bcf86cd799439012",
    "user_id": "507f1f77bcf86cd799439011",
    "snapnotes_id": "507f1f77bcf86cd799439013",
    "prep_pilot_id": "507f1f77bcf86cd799439014",
    "name": "Biology Study Guide",
    "icon": "🧬",
    "color": "#4CAF50",
    "type": "study",
    "google_drive_link": "https://drive.google.com/file/d/1abc123def456",
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  },
  {
    "id": "507f1f77bcf86cd799439015",
    "user_id": "507f1f77bcf86cd799439011",
    "snapnotes_id": null,
    "prep_pilot_id": null,
    "name": "Chemistry Notes",
    "icon": "⚗️",
    "color": "#FF9800",
    "type": "notes",
    "google_drive_link": null,
    "created_at": "2024-01-16T09:20:00Z",
    "updated_at": "2024-01-16T09:20:00Z"
  }
]
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Failed to retrieve notebooks"
}
```

#### Update Notebook
- **PUT** `/api/v1/notebooks/{id}`
- **Authentication:** Required (JWT)
- **Description:** Update a notebook (user can only update their own notebooks)

**Example Request:**
```bash
curl -X PUT http://localhost:8080/api/v1/notebooks/507f1f77bcf86cd799439012 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "Advanced Biology Study Guide",
    "icon": "🔬",
    "color": "#2196F3",
    "type": "research",
    "google_drive_link": "https://drive.google.com/file/d/1updated123"
  }'
```

**Path Parameters:**
- `id` (string): Notebook ID (MongoDB ObjectID)

**Request Body:**
```json
{
  "name": "Advanced Biology Study Guide",
  "icon": "🔬",
  "color": "#2196F3",
  "type": "research",
  "google_drive_link": "https://drive.google.com/file/d/1updated123"
}
```

**Example Response (200 OK):**
```json
{
  "message": "Notebook updated successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request body
```json
{
  "error": "Key: 'Notebook.Name' Error:Tag 'required' validation failed"
}
```
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Failed to update notebook"
}
```

#### Delete Notebook
- **DELETE** `/api/v1/notebooks/{id}`
- **Authentication:** Required (JWT)
- **Description:** Delete a notebook (user can only delete their own notebooks)

**Example Request:**
```bash
curl -X DELETE http://localhost:8080/api/v1/notebooks/507f1f77bcf86cd799439012 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Path Parameters:**
- `id` (string): Notebook ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "message": "Notebook deleted successfully"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `500 Internal Server Error`: Server error
```json
{
  "error": "Failed to delete notebook"
}
```

#### Get Notebook Snapnotes
- **GET** `/api/v1/notebooks/{id}/snapnotes`
- **Authentication:** Required (JWT)
- **Description:** Retrieve snapnotes associated with a specific notebook

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/notebooks/507f1f77bcf86cd799439012/snapnotes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Path Parameters:**
- `id` (string): Notebook ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439013",
  "title": "Biology Study Guide Snapnotes",
  "summaryByChapter": [
    {
      "chapterTitle": "Introduction to Biology",
      "summary": "This chapter covers the fundamental concepts of biology including cell theory, evolution, and the diversity of life. It introduces students to the scientific method and how biological research is conducted.",
      "keyPoints": [
        "Cell theory states that all living things are made of cells",
        "Evolution explains the diversity of life on Earth",
        "DNA contains genetic information",
        "Homeostasis maintains internal balance"
      ]
    },
    {
      "chapterTitle": "Cell Structure and Function",
      "summary": "Detailed exploration of cellular components and their functions in both prokaryotic and eukaryotic cells.",
      "keyPoints": [
        "Nucleus controls cell activities",
        "Mitochondria produce energy",
        "Cell membrane regulates what enters and exits",
        "Ribosomes synthesize proteins"
      ]
    }
  ],
  "flashcards": [
    {
      "chapterTitle": "Introduction to Biology",
      "flashcards": [
        {
          "key term": "Cell",
          "definition": "The basic structural and functional unit of all living organisms"
        },
        {
          "key term": "Homeostasis",
          "definition": "The maintenance of stable internal conditions in an organism"
        }
      ]
    },
    {
      "chapterTitle": "Cell Structure and Function",
      "flashcards": [
        {
          "key term": "Mitochondria",
          "definition": "Organelles that produce ATP energy for cellular processes"
        },
        {
          "key term": "Nucleus",
          "definition": "Control center of the cell containing DNA"
        }
      ]
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Notebook ID is required
```json
{
  "error": "Notebook ID is required"
}
```
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `404 Not Found`: Snapnotes not found or no snapnotes associated with this notebook
```json
{
  "error": "no snapnotes associated with this notebook"
}
```

#### Get Notebook Prep Pilot
- **GET** `/api/v1/notebooks/{id}/prep-pilot`
- **Authentication:** Required (JWT)
- **Description:** Retrieve prep pilot (practice questions) associated with a specific notebook

**Example Request:**
```bash
curl -X GET http://localhost:8080/api/v1/notebooks/507f1f77bcf86cd799439012/prep-pilot \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Path Parameters:**
- `id` (string): Notebook ID (MongoDB ObjectID)

**Example Response (200 OK):**
```json
{
  "id": "507f1f77bcf86cd799439014",
  "notebook_id": "507f1f77bcf86cd799439012",
  "chapters": [
    {
      "chapterTitle": "Introduction to Biology",
      "questions": [
        {
          "question": "What is the basic unit of life?",
          "options": {
            "A": "Atom",
            "B": "Cell",
            "C": "Molecule",
            "D": "Tissue"
          },
          "answer": "B",
          "explanation": "The cell is considered the basic unit of life as it is the smallest unit that can carry out all life processes."
        },
        {
          "question": "Which process explains the diversity of life on Earth?",
          "options": {
            "A": "Photosynthesis",
            "B": "Respiration",
            "C": "Evolution",
            "D": "Digestion"
          },
          "answer": "C",
          "explanation": "Evolution is the process by which species change over time, leading to the diversity of life we see today."
        }
      ]
    },
    {
      "chapterTitle": "Cell Structure and Function",
      "questions": [
        {
          "question": "Which organelle is known as the powerhouse of the cell?",
          "options": {
            "A": "Nucleus",
            "B": "Ribosome",
            "C": "Mitochondria",
            "D": "Endoplasmic Reticulum"
          },
          "answer": "C",
          "explanation": "Mitochondria are called the powerhouse of the cell because they produce ATP, the energy currency of the cell."
        },
        {
          "question": "What controls the activities of the cell?",
          "options": {
            "A": "Cell membrane",
            "B": "Nucleus",
            "C": "Cytoplasm",
            "D": "Vacuole"
          },
          "answer": "B",
          "explanation": "The nucleus contains the cell's DNA and controls all cellular activities including growth, reproduction, and metabolism."
        }
      ]
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Notebook ID is required
```json
{
  "error": "Notebook ID is required"
}
```
- `401 Unauthorized`: Missing or invalid JWT token
```json
{
  "error": "User ID not found in context"
}
```
- `404 Not Found`: Prep pilot not found or no prep pilot associated with this notebook
```json
{
  "error": "no prep pilot associated with this notebook"
}
```

## Error Handling

### Common Error Response Format
All error responses follow this format:
```json
{
  "error": "Error message description"
}
```

### HTTP Status Codes
- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

## Authentication Details

### JWT Token Structure
The JWT token contains the following claims:
- `user_id`: User's MongoDB ObjectID
- `email`: User's email address
- `exp`: Token expiration time (24 hours from issue)

### Token Usage
Include the token in the Authorization header for all protected endpoints:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Environment Variables
The API requires the following environment variables:
- `JWT_SECRET`: Secret key for JWT token signing
- `MONGODB_URI`: MongoDB connection string
- `PORT`: Server port (defaults to 8080)

## Rate Limiting
Currently, no rate limiting is implemented, but it's recommended to implement rate limiting in production environments.

## CORS
CORS configuration should be implemented based on your frontend domain requirements.

## Examples

### Complete User Registration and Notebook Creation Flow

1. **Register a new user:**
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "securepassword123",
    "name": "Jane Student"
  }'
```

2. **Login to get JWT token:**
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "student@example.com",
    "password": "securepassword123"
  }'
```

3. **Create a notebook using the JWT token:**
```bash
curl -X POST http://localhost:8080/api/v1/notebooks/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
  -d '{
    "name": "Biology Study Guide",
    "icon": "🧬",
    "color": "#4CAF50",
    "type": "study"
  }'
```

4. **Get all user notebooks:**
```bash
curl -X GET http://localhost:8080/api/v1/notebooks/user \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## Notes
- All timestamps are in ISO 8601 format
- MongoDB ObjectIDs are represented as hexadecimal strings
- The API uses JSON for all request and response bodies
- Password fields are never returned in API responses for security
- Users can only access and modify their own resources (notebooks, snapnotes, prep pilots)
