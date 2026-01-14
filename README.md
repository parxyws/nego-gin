# Nego-Gin API Route Protection Status

## 1. Authentication & User Management

### Auth

- `POST /api/auth/register` - **Public**
- `POST /api/auth/login` - **Public**
- `POST /api/auth/logout` - **Protected**
- `POST /api/auth/refresh-token` - **Public** (with refresh token)
- `POST /api/auth/forgot-password` - **Public**
- `POST /api/auth/reset-password` - **Public** (with reset token)
- `POST /api/auth/verify-email` - **Public** (with verification token)
- `POST /api/auth/resend-verification` - **Protected**

### Users

- `GET /api/users/me` - **Protected**
- `PUT /api/users/me` - **Protected**
- `PATCH /api/users/me/avatar` - **Protected**
- `DELETE /api/users/me` - **Protected**
- `GET /api/users/:userId` - **Public**
- `GET /api/users/:userId/ratings` - **Public**

### User Addresses

- `GET /api/users/me/addresses` - **Protected**
- `POST /api/users/me/addresses` - **Protected**
- `PUT /api/users/me/addresses/:addressId` - **Protected**
- `DELETE /api/users/me/addresses/:addressId` - **Protected**
- `PATCH /api/users/me/addresses/:addressId/default` - **Protected**

### Roles & Permissions (Admin)

- `GET /api/admin/roles` - **Protected (Admin)**
- `POST /api/admin/roles` - **Protected (Admin)**
- `PUT /api/admin/roles/:roleId` - **Protected (Admin)**
- `DELETE /api/admin/roles/:roleId` - **Protected (Admin)**
- `GET /api/admin/permissions` - **Protected (Admin)**
- `POST /api/admin/users/:userId/roles` - **Protected (Admin)**
- `DELETE /api/admin/users/:userId/roles/:roleId` - **Protected (Admin)**

---

## 2. Product Catalog

### Categories (Marketplace)

- `GET /api/categories` - **Public**
- `GET /api/categories/:categoryId` - **Public**
- `GET /api/categories/:categoryId/products` - **Public**
- `POST /api/admin/categories` - **Protected (Admin)**
- `PUT /api/admin/categories/:categoryId` - **Protected (Admin)**
- `DELETE /api/admin/categories/:categoryId` - **Protected (Admin)**

### Shop Categories (Seller-Specific)

- `GET /api/sellers/shop-categories` - **Protected (Seller)**
- `GET /api/sellers/shop-categories/:shopCategoryId` - **Protected (Seller)**
- `POST /api/sellers/shop-categories` - **Protected (Seller)**
- `PUT /api/sellers/shop-categories/:shopCategoryId` - **Protected (Seller)**
- `DELETE /api/sellers/shop-categories/:shopCategoryId` - **Protected (Seller)**
- `GET /api/shops/:sellerId/categories` - **Public**
- `GET /api/shops/:sellerId/categories/:shopCategoryId` - **Public**

### Tags

- `GET /api/tags` - **Public**
- `GET /api/tags/:tagId/products` - **Public**

### Products

- `GET /api/products` - **Public**
- `GET /api/products/featured` - **Public**
- `GET /api/products/:productId` - **Public**
- `GET /api/products/slug/:slug` - **Public**
- `POST /api/sellers/products` - **Protected (Seller)**
- `PUT /api/sellers/products/:productId` - **Protected (Seller)**
- `DELETE /api/sellers/products/:productId` - **Protected (Seller)**
- `PATCH /api/sellers/products/:productId/status` - **Protected (Seller)**
- `GET /api/sellers/products` - **Protected (Seller)**

### Product Variants

- `GET /api/products/:productId/variants` - **Public**
- `POST /api/sellers/products/:productId/variants` - **Protected (Seller)**
- `PUT /api/sellers/products/:productId/variants/:variantId` - **Protected (Seller)**
- `DELETE /api/sellers/products/:productId/variants/:variantId` - **Protected (Seller)**

### Product Media

- `GET /api/products/:productId/media` - **Public**
- `POST /api/sellers/products/:productId/media` - **Protected (Seller)**
- `PUT /api/sellers/products/:productId/media/:mediaId` - **Protected (Seller)**
- `DELETE /api/sellers/products/:productId/media/:mediaId` - **Protected (Seller)**

---

## 3. Auction System

### Auctions

- `GET /api/auctions` - **Public**
- `GET /api/auctions/live` - **Public**
- `GET /api/auctions/ending-soon` - **Public**
- `GET /api/auctions/:auctionId` - **Public**
- `POST /api/sellers/auctions` - **Protected (Seller)**
- `PUT /api/sellers/auctions/:auctionId` - **Protected (Seller)**
- `DELETE /api/sellers/auctions/:auctionId` - **Protected (Seller)**
- `GET /api/sellers/auctions` - **Protected (Seller)**

### Bids

- `GET /api/auctions/:auctionId/bids` - **Public**
- `POST /api/auctions/:auctionId/bids` - **Protected**
- `GET /api/users/me/bids` - **Protected**
- `GET /api/users/me/bids/active` - **Protected**
- `DELETE /api/auctions/:auctionId/bids/:bidId` - **Protected**

### Watchlists

- `GET /api/users/me/watchlist` - **Protected**
- `POST /api/users/me/watchlist` - **Protected**
- `DELETE /api/users/me/watchlist/:watchlistId` - **Protected**
- `PUT /api/users/me/watchlist/:watchlistId` - **Protected**

---

## 4. Shopping Cart & Checkout

### Cart

- `GET /api/cart` - **Public/Protected** (supports both guest and authenticated)
- `POST /api/cart/items` - **Public/Protected** (supports both guest and authenticated)
- `PUT /api/cart/items/:cartItemId` - **Public/Protected** (supports both guest and authenticated)
- `DELETE /api/cart/items/:cartItemId` - **Public/Protected** (supports both guest and authenticated)
- `DELETE /api/cart` - **Public/Protected** (supports both guest and authenticated)
- `POST /api/cart/merge` - **Protected** (after login)

### Checkout

- `POST /api/checkout/validate` - **Public/Protected** (supports both guest and authenticated)
- `POST /api/checkout` - **Protected** (requires user account)
- `POST /api/checkout/calculate` - **Public/Protected** (supports both guest and authenticated)

---

## 5. Orders

### Orders

- `GET /api/orders` - **Protected**
- `GET /api/orders/:orderId` - **Protected**
- `POST /api/orders/:orderId/cancel` - **Protected**
- `GET /api/sellers/orders` - **Protected (Seller)**
- `PATCH /api/sellers/orders/:orderId/status` - **Protected (Seller)**

### Order Items

- `GET /api/orders/:orderId/items` - **Protected**

---

## 6. Payments

### Payments

- `POST /api/payments` - **Protected**
- `GET /api/payments/:paymentId` - **Protected**
- `POST /api/payments/:paymentId/confirm` - **Protected**
- `POST /api/payments/:paymentId/refund` - **Protected**
- `GET /api/orders/:orderId/payments` - **Protected**
- `POST /api/webhooks/payment/:provider` - **Public** (webhook, verified by signature)

---

## 7. Shipping

### Shipments

- `GET /api/orders/:orderId/shipments` - **Protected**
- `GET /api/shipments/:shipmentId` - **Protected**
- `GET /api/shipments/:shipmentId/tracking` - **Public** (with tracking number)
- `POST /api/sellers/orders/:orderId/shipments` - **Protected (Seller)**
- `PUT /api/sellers/shipments/:shipmentId` - **Protected (Seller)**

---

## 8. Reviews & Ratings

### Reviews

- `GET /api/products/:productId/reviews` - **Public**
- `GET /api/sellers/:sellerId/reviews` - **Public**
- `POST /api/reviews` - **Protected**
- `PUT /api/reviews/:reviewId` - **Protected**
- `DELETE /api/reviews/:reviewId` - **Protected**
- `GET /api/users/me/reviews` - **Protected**
- `POST /api/reviews/:reviewId/media` - **Protected**

### Review Moderation (Admin)

- `GET /api/admin/reviews/pending` - **Protected (Admin)**
- `PATCH /api/admin/reviews/:reviewId/approve` - **Protected (Admin)**
- `PATCH /api/admin/reviews/:reviewId/reject` - **Protected (Admin)**

---

## 9. Notifications

### Notifications

- `GET /api/notifications` - **Protected**
- `GET /api/notifications/unread` - **Protected**
- `PATCH /api/notifications/:notificationId/read` - **Protected**
- `PATCH /api/notifications/read-all` - **Protected**
- `DELETE /api/notifications/:notificationId` - **Protected**

---

## 10. Search & Discovery

### Search

- `GET /api/search` - **Public**
- `GET /api/search/suggestions` - **Public**
- `GET /api/search/filters` - **Public**

---

## 11. Analytics & Reports (Seller/Admin)

### Seller Dashboard

- `GET /api/sellers/dashboard` - **Protected (Seller)**
- `GET /api/sellers/analytics/sales` - **Protected (Seller)**
- `GET /api/sellers/analytics/products` - **Protected (Seller)**
- `GET /api/sellers/analytics/auctions` - **Protected (Seller)**

### Admin Reports

- `GET /api/admin/reports/sales` - **Protected (Admin)**
- `GET /api/admin/reports/users` - **Protected (Admin)**
- `GET /api/admin/reports/auctions` - **Protected (Admin)**
- `GET /api/admin/reports/revenue` - **Protected (Admin)**

---

## 12. WebSocket/Real-time Endpoints

### Real-time Features

- `WS /api/ws/auctions/:auctionId` - **Public** (but authenticated users get enhanced features)
- `WS /api/ws/notifications` - **Protected**

---

## Legend

- **Public** - No authentication required
- **Protected** - Requires authentication (logged-in user)
- **Protected (Seller)** - Requires authentication + seller role
- **Protected (Admin)** - Requires authentication + admin role
- **Public/Protected** - Works for both guest and authenticated users (typically cart/checkout)