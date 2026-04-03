#  Finance Data Processing & Access Control Backend

##  Project Overview

This project is a backend system developed as part of an assignment for **Zorvyn FinTech Pvt. Ltd.**.

The goal is to build a scalable backend for a **Finance Dashboard System** with:

* Financial record management
* Role-based access control (RBAC)
* Dashboard analytics
* Clean API design

---

##  Author

**Keshav Raj**
 [contacttokeshavraj@gmail.com](mailto:contacttokeshavraj@gmail.com)
 https://www.linkedin.com/in/keshavraj18

---

##  Tech Stack

* Golang (Go)
* Fiber Framework
* MongoDB
* JWT Authentication

---

##  Project Structure

finance-dashboard-backend/
├── cmd/
├── config/
├── controller/
├── repository/
├── models/
├── middleware/
├── utils/
├── db/
├── router/
├── constants/

---

##  Roles

* Viewer → View dashboard
* Analyst → View + Create/Update records
* Admin → Full access

---

##  Features

* User Management (Register/Login)
* Role-based Access Control
* Record CRUD + Filtering
* Dashboard Analytics
* Pagination
* Validation & Error Handling

---

##  Setup

```bash
go mod tidy
go run cmd/main.go
```

---

##  APIs

### Auth

POST /api/auth/register
POST /api/auth/login

### Users

GET /api/users
PATCH /api/users/:id/role
PATCH /api/users/:id/status

### Records

POST /api/records
GET /api/records?page=1&limit=10
PUT /api/records/:id
DELETE /api/records/:id

### Dashboard

GET /api/dashboard/summary
GET /api/dashboard/category
GET /api/dashboard/monthly
GET /api/dashboard/recent

---

##  Conclusion

This project demonstrates strong backend fundamentals including:

* Clean architecture
* RBAC implementation
* Aggregation logic
* Scalable API design

---
