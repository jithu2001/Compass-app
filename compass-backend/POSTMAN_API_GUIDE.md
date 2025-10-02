# Compass Backend API - Postman Guide

This document provides detailed Postman examples for all API endpoints in the Compass Backend project.

## Base Configuration

### Environment Variables
Set up these variables in your Postman environment:

```
base_url: http://localhost:8080
access_token: {{access_token}}
refresh_token: {{refresh_token}}
```

### Headers for Protected Endpoints
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

---

## Authentication Endpoints

### 1. Sign In
**Method:** `POST`  
**URL:** `{{base_url}}/auth/signin`  
**Headers:**
```
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "email": "admin@compass.com",
    "password": "compass_password"
}
```

**Test Script (to save tokens):**
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.access_token);
    pm.environment.set("refresh_token", jsonData.refresh_token);
    pm.environment.set("user_id", jsonData.user.user_id);
}
```

**Expected Response:**
```json
{
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "user_id": 1,
        "email": "admin@compass.com",
        "full_name": "Admin User",
        "role": "admin"
    }
}
```

### 2. Refresh Token
**Method:** `POST`  
**URL:** `{{base_url}}/auth/refresh`  
**Headers:**
```
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "refresh_token": "{{refresh_token}}"
}
```

**Test Script:**
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.access_token);
    pm.environment.set("refresh_token", jsonData.refresh_token);
}
```

### 3. Logout
**Method:** `POST`  
**URL:** `{{base_url}}/api/auth/logout`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

---

## User Management (Admin Only)

### 4. Create User (Admin Only)
**Method:** `POST`  
**URL:** `{{base_url}}/api/users`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "full_name": "John Doe",
    "email": "john.doe@compass.com",
    "password": "securePassword123",
    "role": "user"
}
```

**Alternative - Invite User (no password):**
```json
{
    "full_name": "Jane Smith",
    "email": "jane.smith@compass.com",
    "role": "admin"
}
```

### 5. List All Users (Admin Only)
**Method:** `GET`  
**URL:** `{{base_url}}/api/users`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

### 6. Update User Status (Admin Only)
**Method:** `PATCH`  
**URL:** `{{base_url}}/api/users/2/status`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "status": "active"
}
```

**Valid status values:** `pending`, `active`, `disabled`

---

## Project Management

### 7. Create Project with Full Details
**Method:** `POST`  
**URL:** `{{base_url}}/api/projects`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON) - Complete Example:**
```json
{
    "project_name": "City Tower Windows Replacement",
    "company_name": "ABC Construction Ltd",
    "company_address": "123 Main Street, London, UK",
    "project_type": "windows",
    "specifications": [
        {
            "version_no": 1,
            "colour": "RAL 7016 Anthracite Grey",
            "ironmongery": "Stainless Steel Handles",
            "u_value": 1.4,
            "g_value": 0.7,
            "vents": "Top hung with restrictors",
            "acoustics": "45dB reduction",
            "sbd": true,
            "pas24": true,
            "restrictors": true,
            "special_comments": "All windows must be triple glazed with laminated inner pane",
            "attachment_url": "https://example.com/specs/city-tower-v1.pdf"
        }
    ],
    "rfis": [
        {
            "question_text": "Is fire rating required for windows above 18m?"
        },
        {
            "question_text": "Are automatic closers needed for ground floor windows?"
        }
    ]
}
```

**Body (raw JSON) - Minimal Example:**
```json
{
    "project_name": "Small Office Renovation",
    "project_type": "doors"
}
```

### 8. List All Projects
**Method:** `GET`  
**URL:** `{{base_url}}/api/projects`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

### 9. Get Project Details
**Method:** `GET`  
**URL:** `{{base_url}}/api/projects/1`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

### 10. Update Project Status
**Method:** `PATCH`  
**URL:** `{{base_url}}/api/projects/1/status`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "status": "progress"
}
```

**Valid status values:** `not_yet_started`, `progress`, `completed`

### 11. Delete Project (Admin Only)
**Method:** `DELETE`  
**URL:** `{{base_url}}/api/projects/1`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

---

## Project Specifications

### 12. Add New Specification to Project
**Method:** `POST`  
**URL:** `{{base_url}}/api/projects/1/specifications`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON) - Version 2:**
```json
{
    "colour": "RAL 9010 Pure White",
    "ironmongery": "Chrome Handles with Lock",
    "u_value": 1.2,
    "g_value": 0.65,
    "vents": "Side hung with micro ventilation",
    "acoustics": "50dB reduction",
    "sbd": true,
    "pas24": true,
    "restrictors": false,
    "special_comments": "Updated specification after client review. Changed color from grey to white.",
    "attachment_url": "https://example.com/specs/city-tower-v2.pdf"
}
```

### 13. Get All Specifications for a Project
**Method:** `GET`  
**URL:** `{{base_url}}/api/projects/1/specifications`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

---

## Project RFIs (Requests for Information)

### 14. Create New RFI
**Method:** `POST`  
**URL:** `{{base_url}}/api/projects/1/rfis`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "question_text": "What is the required lead time for window delivery?"
}
```

### 15. Get All RFIs for a Project
**Method:** `GET`  
**URL:** `{{base_url}}/api/projects/1/rfis`  
**Headers:**
```
Authorization: Bearer {{access_token}}
```

### 16. Answer an RFI
**Method:** `PATCH`  
**URL:** `{{base_url}}/api/rfis/1/answer`  
**Headers:**
```
Authorization: Bearer {{access_token}}
Content-Type: application/json
```

**Body (raw JSON):**
```json
{
    "answer_value": "yes"
}
```

**Valid answer values:** `yes`, `no`

---

## Health Check

### 17. Service Health Check
**Method:** `GET`  
**URL:** `{{base_url}}/health`  
**Headers:** None required

**Expected Response:**
```json
{
    "status": "healthy",
    "service": "compass-backend"
}
```

---

## Postman Collection Organization

### Recommended Folder Structure:
```
Compass Backend API
├── Authentication
│   ├── Sign In
│   ├── Refresh Token
│   └── Logout
├── User Management (Admin)
│   ├── Create User
│   ├── List Users
│   └── Update User Status
├── Projects
│   ├── Create Project
│   ├── List Projects
│   ├── Get Project
│   ├── Update Status
│   └── Delete Project
├── Specifications
│   ├── Add Specification
│   └── List Specifications
├── RFIs
│   ├── Create RFI
│   ├── List RFIs
│   └── Answer RFI
└── Health
    └── Health Check
```

### Pre-request Script (Collection Level):
```javascript
// Add timestamp to requests for logging
pm.request.headers.add({
    key: 'X-Request-Time',
    value: new Date().toISOString()
});
```

### Test Script (Collection Level):
```javascript
// Log response time
console.log(`Response Time: ${pm.response.responseTime}ms`);

// Check for common errors
if (pm.response.code === 401) {
    console.log("Unauthorized - Token may have expired");
}

if (pm.response.code === 403) {
    console.log("Forbidden - Insufficient permissions");
}

if (pm.response.code === 500) {
    console.log("Server Error:", pm.response.json());
}
```

---

## Testing Workflow

### Initial Setup:
1. Import this collection into Postman
2. Create an environment with `base_url` set to your server
3. Run the Sign In request first to get tokens

### Typical Testing Flow:
1. **Sign In** - Get access token
2. **Create Project** - Create a new project with specifications
3. **List Projects** - Verify project was created
4. **Add RFI** - Add questions to the project
5. **Answer RFI** - Provide answers
6. **Update Project Status** - Mark as in progress
7. **Add New Specification** - Create version 2

### Admin Testing:
1. Sign in with admin credentials
2. Create new users
3. List all users
4. Update user status
5. Test admin-only endpoints

### Error Testing:
- Try accessing endpoints without tokens (expect 401)
- Try admin endpoints with regular user (expect 403)
- Try invalid data formats (expect 400)
- Try non-existent resources (expect 404)

---

## Notes

- All timestamps are in UTC format
- Tokens expire after the duration set in the server config
- File uploads for attachments should be handled separately
- The `version_no` for specifications auto-increments per project
- Users invited without passwords receive email invitations
- Project deletion is a soft delete (marked as deleted, not removed)