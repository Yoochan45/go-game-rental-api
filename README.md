# Video Game Rental API

## Overview
Video Game Rental API adalah sistem backend berbasis Golang (Echo Framework) untuk platform penyewaan game fisik seperti kaset dan console.  
Proyek ini menerapkan sistem multi-role (Super Admin, Admin, Partner, Customer), serta fitur KYC, sistem pembayaran, review, dan approval flow.

---

## Tech Stack
| Layer | Teknologi |
|-------|------------|
| Backend | Go (Echo v4) |
| Database | PostgreSQL / Supabase |
| ORM / Query | GORM |
| Authentication | JWT + Refresh Token |
| File Storage | Supabase Storage |
| Payment Gateway | Stripe / Midtrans |
| Validation | go-playground/validator v10 |
| Logging | logrus |
| Documentation | Swagger (swaggo) |
| CI/CD | Heroku / Railway |
| Testing | testify, mockgen |

---

## Modules & Features

### Auth
- Register & Login (default role: `customer`)
- Refresh Token
- Role-based Access Control (RBAC)
- JWT Middleware

### User
- View & Edit Profile
- KYC Upload
- Deposit / Wallet Management

### Partner
- Apply for Partner (KYC + Verification)
- Manage Own Listings (CRUD Game)
- View Bookings & Payments related to their listings

### Admin
- Approve or Reject KYC Applications
- Approve or Reject Listings
- Manage Users
- Handle Disputes
- View Reports

### Super Admin
- Full system ownership
- Manage Admins
- Manage global settings
- Emergency access

### Catalog
- List Games
- Filter by Category, Availability, or Price
- Manage Stock & Status

### Booking & Payment
- Create Booking (pending → paid → confirmed → completed)
- Integrate with Payment Gateway (Stripe / Midtrans)
- Webhook Handling for Payment Confirmation
- Refund & Cancellation Logic

### Review
- CRUD Review untuk rental yang sudah selesai

### Reporting & Audit
- Log every admin action (admin_logs)
- Generate reports for financial and operational analysis

---

## Detailed Business Flow

### Partner Onboarding Flow
1. User register → role: `customer`
2. Customer uploads KYC documents via `/users/kyc/upload`
3. Customer submits partner application via `/partner/apply`
4. Admin reviews KYC via `/admin/kyc/:id/approve`
5. If approved → user.role = `partner`
6. Partner can now create game listings

### Game Listing Flow
1. Partner creates game listing via `/partner/games`
2. Listing status = `pending_approval`
3. Admin reviews listing via `/admin/listings/:id/approve`
4. If approved → game.is_active = true, visible to customers

### Booking & Rental Flow
1. Customer browses approved games via `/games`
2. Customer creates booking via `/bookings` → status: `pending_payment`
3. Customer pays via `/payments/:booking_id/pay`
4. Payment webhook confirms → booking status: `confirmed`
5. **Partner confirms item handover** → status: `active`
6. Customer returns item → Partner confirms return → status: `completed`
7. Customer can leave review via `/bookings/:id/review`

### Dispute & Refund Flow
1. Customer/Partner reports issue via `/disputes/create`
2. Admin investigates via `/admin/disputes/:id`
3. Admin decides refund/resolution via `/admin/disputes/:id/resolve`

---

## Entity Relationship Diagram (ERD) - Summary
- users (id, email, password_hash, role, deposit, is_active, created_at, updated_at)
- roles (id, name, permissions)
- partner_applications (id, user_id, status, kyc_doc_id, submitted_at, decided_by, decided_at, rejection_reason)
- games (id, partner_id, name, category_id, stock, rental_price, description, is_active, approved_by_admin, approval_status, created_at)
- categories (id, name, description)
- bookings (id, user_id, game_id, partner_id, start_at, end_at, total_price, status, created_at, updated_at)
- payments (id, booking_id, provider, provider_payment_id, amount, status, paid_at, refunded_at)
- reviews (id, booking_id, user_id, game_id, rating, comment, created_at)
- rental_history (id, booking_id, user_id, game_snapshot, price, start_at, end_at, returned_at)
- wallets (id, user_id, balance, created_at, updated_at)
- wallet_transactions (id, wallet_id, type, amount, description, reference_id, created_at)
- kyc_documents (id, user_id, doc_type, url, status, uploaded_at, verified_at)
- disputes (id, booking_id, reporter_id, type, description, status, resolution, created_at, resolved_at)
- admin_logs (id, admin_id, action, target_type, target_id, metadata, created_at)
- system_settings (id, key, value, description, updated_by, updated_at)

---

## API Endpoint Pattern

| Resource | Method | Endpoint | Description |
|-----------|---------|----------|-------------|
| **Auth** | POST | /auth/register | Register user |
|  | POST | /auth/login | Login |
|  | POST | /auth/refresh | Refresh token |
| **Users** | GET | /users/me | Get current user profile |
|  | PUT | /users/me | Update profile |
| **Partner** | POST | /partner/apply | Submit partner application (KYC) |
|  | GET | /partner/applications | Get all partner applications *(admin only)* |
|  | PATCH | /partner/applications/:id/approve | Approve or reject partner *(admin)* |
| **Games** | GET | /games | Get all games |
|  | GET | /games/:id | Get game detail |
|  | POST | /partner/games | Create new game listing *(partner)* |
|  | PATCH | /partner/games/:id | Update own game listing *(partner)* |
| **Bookings** | POST | /bookings | Create booking *(customer)* |
|  | GET | /bookings/:id | Get booking detail *(authorized only)* |
|  | PATCH | /bookings/:id/cancel | Cancel booking *(customer)* |
| **Payments** | POST | /payments/:booking_id/pay | Make payment |
|  | POST | /webhooks/payments | Handle payment webhook *(system)* |
| **Reviews** | POST | /bookings/:id/review | Add review after completed booking |
| **Admin** | GET | /admin/users | View all users |
|  | PATCH | /admin/users/:id/ban | Ban / unban user |
|  | GET | /admin/kyc | View pending KYC submissions |
|  | PATCH | /admin/kyc/:id/approve | Approve / reject KYC |
|  | GET | /admin/listings | View pending listings |
|  | PATCH | /admin/listings/:id/approve | Approve / reject listing |
|  | GET | /admin/disputes | Handle dispute cases |
|  | GET | /admin/reports | View financial / system reports |
| **Superadmin** | GET | /superadmin/admins | View all admins |
|  | POST | /superadmin/admins | Create new admin |
|  | DELETE | /superadmin/admins/:id | Remove admin |
|  | GET | /superadmin/system/logs | View system logs & audit trail |
|  | PATCH | /superadmin/system/settings | Update global system settings |
|  | POST | /superadmin/system/emergency | Trigger emergency system recovery |
| **KYC** | POST | /users/kyc/upload | Upload KYC documents |
|  | GET | /users/kyc/status | Check KYC status |
| **Wallet** | GET | /wallet/balance | Get wallet balance |
|  | POST | /wallet/deposit | Add funds to wallet |
|  | GET | /wallet/transactions | Get transaction history |
| **Disputes** | POST | /disputes/create | Report dispute |
|  | GET | /disputes/my | Get user's disputes |
| **Partner Dashboard** | GET | /partner/dashboard | Partner analytics |
|  | GET | /partner/bookings | View bookings for partner's games |
|  | PATCH | /partner/bookings/:id/confirm-handover | Confirm item handover |
|  | PATCH | /partner/bookings/:id/confirm-return | Confirm item return |

---

## Security & Authentication

### Public Endpoints (No authentication required)
- `POST /auth/register` - User registration
- `POST /auth/login` - User login  
- `POST /auth/refresh` - Refresh JWT token
- `GET /games` - Browse game catalog
- `GET /games/:id` - View game details
- `POST /webhooks/payments` - Payment gateway webhooks

### Protected Endpoints
All other endpoints require valid JWT token in Authorization header:
```
Authorization: Bearer <jwt_token>
```

### Role-Based Access Control (RBAC)
- **Customer**: Can book games, manage profile, wallet operations
- **Partner**: Customer permissions + manage game listings, view bookings
- **Admin**: Partner permissions + approve KYC/listings, handle disputes
- **Super Admin**: Full system access + manage admins

---

## Status Definitions

### User Roles
- `customer` - Default role, can book games
- `partner` - Can list games for rental
- `admin` - Can approve/reject applications and listings
- `super_admin` - Full system access

### Booking Status
- `pending_payment` - Awaiting payment
- `confirmed` - Payment received
- `active` - Item handed over, rental in progress
- `completed` - Item returned successfully
- `cancelled` - Booking cancelled
- `disputed` - Under dispute resolution

### Payment Status
- `pending` - Payment initiated
- `paid` - Payment successful
- `failed` - Payment failed
- `refunded` - Payment refunded

### KYC Status
- `pending` - Documents uploaded, awaiting review
- `approved` - KYC verified
- `rejected` - KYC rejected
- `resubmission_required` - Need additional documents

---

## Third-Party Integration
- **Database & Storage**: Supabase (Postgres + Storage)
- **Payment Gateway**: Stripe / Midtrans (sandbox mode)
- **Email Notification**: SendGrid / Mailgun
- **Error Tracking**: Sentry
- **Deployment**: Heroku
- **Docs**: Swagger auto-generated

---

## Setup Guide
1. Clone repository
   ```bash
   git clone https://github.com/Yoochan45/go-game-rental-api.git
   cd go-game-rental-api
   ```

2. Install dependencies
   ```bash
   go mod tidy
   ```

3. Setup environment variables
   ```bash
   cp .env.example .env
   # Edit .env dengan konfigurasi Anda
   ```

4. Run database migrations
   ```bash
   go run main.go migrate
   ```

5. Seed initial data (optional)
   ```bash
   go run main.go seed
   ```

6. Run the application
   ```bash
   go run main.go
   # atau
   make run
   ```

---
## Constributor
Aisiya Qutwatunnada (Yoochan45)