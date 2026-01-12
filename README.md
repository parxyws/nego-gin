# Nego-Gin API Documentation

## 1. Authentication & User Management

### Auth
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh-token` - Refresh access token
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password with token
- `POST /api/auth/verify-email` - Verify email with token
- `POST /api/auth/resend-verification` - Resend verification email

### Users
- `GET /api/users/me` - Get current user profile
- `PUT /api/users/me` - Update current user profile
- `PATCH /api/users/me/avatar` - Update avatar
- `DELETE /api/users/me` - Delete account (soft delete)
- `GET /api/users/:userId` - Get public user profile
- `GET /api/users/:userId/ratings` - Get user seller ratings

### User Addresses
- `GET /api/users/me/addresses` - List user addresses
- `POST /api/users/me/addresses` - Create address
- `PUT /api/users/me/addresses/:addressId` - Update address
- `DELETE /api/users/me/addresses/:addressId` - Delete address
- `PATCH /api/users/me/addresses/:addressId/default` - Set default address

### Roles & Permissions (Admin)
- `GET /api/admin/roles` - List roles
- `POST /api/admin/roles` - Create role
- `PUT /api/admin/roles/:roleId` - Update role
- `DELETE /api/admin/roles/:roleId` - Delete role
- `GET /api/admin/permissions` - List permissions
- `POST /api/admin/users/:userId/roles` - Assign role to user
- `DELETE /api/admin/users/:userId/roles/:roleId` - Remove role from user

---

## 2. Product Catalog

### Categories (Marketplace)
- `GET /api/categories` - List all categories (tree structure)
- `GET /api/categories/:categoryId` - Get category details
- `GET /api/categories/:categoryId/products` - List products in category
- `POST /api/admin/categories` - Create category (admin)
- `PUT /api/admin/categories/:categoryId` - Update category (admin)
- `DELETE /api/admin/categories/:categoryId` - Delete category (admin)

### Shop Categories (Seller-Specific)
- `GET /api/sellers/shop-categories` - List seller's shop categories
- `GET /api/sellers/shop-categories/:shopCategoryId` - Get shop category details
- `POST /api/sellers/shop-categories` - Create shop category
- `PUT /api/sellers/shop-categories/:shopCategoryId` - Update shop category
- `DELETE /api/sellers/shop-categories/:shopCategoryId` - Delete shop category
- `GET /api/shops/:sellerId/categories` - List public categories for a specific shop
- `GET /api/shops/:sellerId/categories/:shopCategoryId` - Get public shop category details

### Tags
- `GET /api/tags` - List all tags
- `GET /api/tags/:tagId/products` - Get products by tag

### Products
- `GET /api/products` - List products (with filters, search, pagination)
- `GET /api/products/featured` - List featured products
- `GET /api/products/:productId` - Get product details
- `GET /api/products/slug/:slug` - Get product by slug
- `POST /api/sellers/products` - Create product (seller)
- `PUT /api/sellers/products/:productId` - Update product (seller)
- `DELETE /api/sellers/products/:productId` - Delete product (seller)
- `PATCH /api/sellers/products/:productId/status` - Update product status
- `GET /api/sellers/products` - List seller's products

### Product Variants
- `GET /api/products/:productId/variants` - List product variants
- `POST /api/sellers/products/:productId/variants` - Create variant
- `PUT /api/sellers/products/:productId/variants/:variantId` - Update variant
- `DELETE /api/sellers/products/:productId/variants/:variantId` - Delete variant

### Product Media
- `GET /api/products/:productId/media` - List product media
- `POST /api/sellers/products/:productId/media` - Upload media
- `PUT /api/sellers/products/:productId/media/:mediaId` - Update media order
- `DELETE /api/sellers/products/:productId/media/:mediaId` - Delete media

---

## 3. Auction System

### Auctions
- `GET /api/auctions` - List active auctions (with filters)
- `GET /api/auctions/live` - List live auctions
- `GET /api/auctions/ending-soon` - List auctions ending soon
- `GET /api/auctions/:auctionId` - Get auction details
- `POST /api/sellers/auctions` - Create auction (seller)
- `PUT /api/sellers/auctions/:auctionId` - Update auction (seller)
- `DELETE /api/sellers/auctions/:auctionId` - Cancel auction (seller)
- `GET /api/sellers/auctions` - List seller's auctions

### Bids
- `GET /api/auctions/:auctionId/bids` - Get auction bid history
- `POST /api/auctions/:auctionId/bids` - Place bid
- `GET /api/users/me/bids` - Get user's bid history
- `GET /api/users/me/bids/active` - Get user's active bids
- `DELETE /api/auctions/:auctionId/bids/:bidId` - Cancel bid (if allowed)

### Watchlists
- `GET /api/users/me/watchlist` - Get user's watchlist
- `POST /api/users/me/watchlist` - Add product/auction to watchlist
- `DELETE /api/users/me/watchlist/:watchlistId` - Remove from watchlist
- `PUT /api/users/me/watchlist/:watchlistId` - Update notification preferences

---

## 4. Shopping Cart & Checkout

### Cart
- `GET /api/cart` - Get cart (user or session)
- `POST /api/cart/items` - Add item to cart
- `PUT /api/cart/items/:cartItemId` - Update cart item quantity
- `DELETE /api/cart/items/:cartItemId` - Remove item from cart
- `DELETE /api/cart` - Clear cart
- `POST /api/cart/merge` - Merge guest cart with user cart (after login)

### Checkout
- `POST /api/checkout/validate` - Validate cart before checkout
- `POST /api/checkout` - Create order from cart
- `POST /api/checkout/calculate` - Calculate totals (tax, shipping)

---

## 5. Orders

### Orders
- `GET /api/orders` - List user's orders
- `GET /api/orders/:orderId` - Get order details
- `POST /api/orders/:orderId/cancel` - Cancel order
- `GET /api/sellers/orders` - List seller's orders
- `PATCH /api/sellers/orders/:orderId/status` - Update order status (seller)

### Order Items
- `GET /api/orders/:orderId/items` - List order items

---

## 6. Payments

### Payments
- `POST /api/payments` - Create payment
- `GET /api/payments/:paymentId` - Get payment details
- `POST /api/payments/:paymentId/confirm` - Confirm payment
- `POST /api/payments/:paymentId/refund` - Request refund
- `GET /api/orders/:orderId/payments` - Get order payments
- `POST /api/webhooks/payment/:provider` - Payment provider webhooks

---

## 7. Shipping

### Shipments
- `GET /api/orders/:orderId/shipments` - Get order shipments
- `GET /api/shipments/:shipmentId` - Get shipment details
- `GET /api/shipments/:shipmentId/tracking` - Get tracking history
- `POST /api/sellers/orders/:orderId/shipments` - Create shipment (seller)
- `PUT /api/sellers/shipments/:shipmentId` - Update shipment (seller)

---

## 8. Reviews & Ratings

### Reviews
- `GET /api/products/:productId/reviews` - List product reviews
- `GET /api/sellers/:sellerId/reviews` - List seller reviews
- `POST /api/reviews` - Create review (for order item)
- `PUT /api/reviews/:reviewId` - Update review
- `DELETE /api/reviews/:reviewId` - Delete review
- `GET /api/users/me/reviews` - Get user's reviews
- `POST /api/reviews/:reviewId/media` - Upload review media

### Review Moderation (Admin)
- `GET /api/admin/reviews/pending` - List pending reviews
- `PATCH /api/admin/reviews/:reviewId/approve` - Approve review
- `PATCH /api/admin/reviews/:reviewId/reject` - Reject review

---

## 9. Notifications

### Notifications
- `GET /api/notifications` - List user notifications
- `GET /api/notifications/unread` - Get unread notifications
- `PATCH /api/notifications/:notificationId/read` - Mark as read
- `PATCH /api/notifications/read-all` - Mark all as read
- `DELETE /api/notifications/:notificationId` - Delete notification

---

## 10. Search & Discovery

### Search
- `GET /api/search` - Global search (products, auctions)
- `GET /api/search/suggestions` - Search autocomplete
- `GET /api/search/filters` - Get available filters

---

## 11. Analytics & Reports (Seller/Admin)

### Seller Dashboard
- `GET /api/sellers/dashboard` - Get dashboard stats
- `GET /api/sellers/analytics/sales` - Sales analytics
- `GET /api/sellers/analytics/products` - Product performance
- `GET /api/sellers/analytics/auctions` - Auction performance

### Admin Reports
- `GET /api/admin/reports/sales` - Sales reports
- `GET /api/admin/reports/users` - User statistics
- `GET /api/admin/reports/auctions` - Auction statistics
- `GET /api/admin/reports/revenue` - Revenue reports

---

## 12. WebSocket/Real-time Endpoints

### Real-time Features
- `WS /api/ws/auctions/:auctionId` - Real-time auction updates
- `WS /api/ws/notifications` - Real-time notifications