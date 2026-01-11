-- ============================================================================
-- USERS & AUTHENTICATION (ULID)
-- ============================================================================

CREATE TABLE roles
(
    role_id     SERIAL PRIMARY KEY,
    role_name   VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMP NULL
);

CREATE TABLE permissions
(
    permission_id SERIAL PRIMARY KEY,
    name          VARCHAR(200) NOT NULL UNIQUE,
    description   TEXT,
    resource      VARCHAR(100) NOT NULL,
    action        VARCHAR(50)  NOT NULL,
    created_at    TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP NULL
);

CREATE TABLE users
(
    user_id            CHAR(26) PRIMARY KEY,
    email              VARCHAR(255) NOT NULL UNIQUE,
    password_hash      VARCHAR(255) NOT NULL,
    username           VARCHAR(100) NOT NULL UNIQUE,
    first_name         VARCHAR(100),
    last_name          VARCHAR(100),
    phone              VARCHAR(20),
    avatar_url         TEXT,
    account_status     VARCHAR(20)  NOT NULL DEFAULT 'active', -- ('active', 'suspended', 'banned', 'pending')
    is_verified        TIMESTAMP    NULL,
    last_login         TIMESTAMP,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMP NULL
);

CREATE TABLE role_permissions
(
    role_id       SERIAL NOT NULL,
    permission_id SERIAL NOT NULL,
    assigned_at   TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
    assigned_by   CHAR(26),
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles (role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions (permission_id) ON DELETE CASCADE,
    FOREIGN KEY (assigned_by) REFERENCES users (user_id) ON DELETE SET NULL
);

CREATE TABLE user_roles
(
    user_role_id SERIAL PRIMARY KEY,
    user_id      CHAR(26)  NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    role_id      INTEGER   NOT NULL REFERENCES roles (role_id) ON DELETE CASCADE,
    assigned_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, role_id)
);

CREATE TABLE user_addresses
(
    address_id     SERIAL PRIMARY KEY,
    user_id        CHAR(26)     NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    address_type   VARCHAR(20)  NOT NULL, -- ('billing', 'shipping', 'both')
    is_default     BOOLEAN      NOT NULL DEFAULT FALSE,
    recipient_name VARCHAR(200) NOT NULL,
    address_line1  VARCHAR(255) NOT NULL,
    address_line2  VARCHAR(255),
    city           VARCHAR(100) NOT NULL,
    state_province VARCHAR(100),
    postal_code    VARCHAR(20)  NOT NULL,
    country_code   CHAR(2)      NOT NULL,
    phone          VARCHAR(20),
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP NULL
);

-- ============================================================================
-- PRODUCT CATALOG (Categories: SERIAL, Products: ULID)
-- ============================================================================

CREATE TABLE categories
(
    category_id        SERIAL PRIMARY KEY,
    parent_category_id INTEGER      REFERENCES categories (category_id) ON DELETE SET NULL,
    category_name      VARCHAR(100) NOT NULL,
    slug               VARCHAR(100) NOT NULL UNIQUE,
    description        TEXT,
    image_url          TEXT,
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE,
    sort_order         INTEGER               DEFAULT 0,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMP NULL
);

CREATE TABLE tags
(
    tag_id     SERIAL PRIMARY KEY,
    tag_name   VARCHAR(50) NOT NULL UNIQUE,
    slug       VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE TABLE products
(
    product_id          CHAR(26) PRIMARY KEY,
    seller_id           CHAR(26)     NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT,
    category_id         INTEGER      NOT NULL REFERENCES categories (category_id) ON DELETE RESTRICT,
    product_type        VARCHAR(20)  NOT NULL CHECK (product_type IN ('standard', 'auction')),
    sku                 VARCHAR(100) NOT NULL UNIQUE,
    name                VARCHAR(255) NOT NULL,
    slug                VARCHAR(255) NOT NULL UNIQUE,
    description         TEXT,
    short_description   TEXT,

    -- Pricing (for standard products)
    base_price          DECIMAL(10, 2),
    sale_price          DECIMAL(10, 2),
    cost_price          DECIMAL(10, 2),

    -- Inventory
    stock_quantity      INTEGER               DEFAULT 0,
    low_stock_threshold INTEGER               DEFAULT 5,
    is_unlimited_stock  BOOLEAN      NOT NULL DEFAULT FALSE,

    -- Status
    status              VARCHAR(20)  NOT NULL DEFAULT 'draft', -- ('draft', 'active', 'inactive', 'archived')
    is_featured         BOOLEAN      NOT NULL DEFAULT FALSE,

    -- Metadata
    weight_kg           DECIMAL(8, 2),
    dimensions_cm       VARCHAR(50),                           -- Format: "LxWxH"
    brand               VARCHAR(100),
    condition           VARCHAR(20),

    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL,

    -- Constraints
    CONSTRAINT price_check CHECK (
        (product_type = 'standard' AND base_price IS NOT NULL) OR
        (product_type = 'auction')
        )
);

CREATE TABLE product_variants
(
    variant_id       SERIAL PRIMARY KEY,
    product_id       CHAR(26)     NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    sku              VARCHAR(100) NOT NULL UNIQUE,
    variant_name     VARCHAR(255) NOT NULL,

    -- Variant-specific pricing
    price_adjustment DECIMAL(10, 2)        DEFAULT 0,

    -- Variant attributes (JSON for flexibility: color, size, etc.)
    attributes       JSONB        NOT NULL,

    -- Inventory
    stock_quantity   INTEGER               DEFAULT 0,

    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at       TIMESTAMP NULL
);

CREATE TABLE product_attributes
(
    attribute_id    SERIAL PRIMARY KEY,
    product_id      CHAR(26)     NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    attribute_name  VARCHAR(100) NOT NULL,
    attribute_value TEXT         NOT NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (product_id, attribute_name)
);

CREATE TABLE product_tags
(
    product_id CHAR(26)  NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    tag_id     INTEGER   NOT NULL REFERENCES tags (tag_id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (product_id, tag_id)
);

CREATE TABLE product_media
(
    media_id      SERIAL PRIMARY KEY,
    product_id    CHAR(26)    NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    media_type    VARCHAR(20) NOT NULL CHECK (media_type IN ('image', 'video')),
    media_url     TEXT        NOT NULL,
    thumbnail_url TEXT,
    alt_text      VARCHAR(255),
    sort_order    INTEGER              DEFAULT 0,
    is_primary    BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- AUCTION SYSTEM (ULID for time-ordered nature)
-- ============================================================================

CREATE TABLE auctions
(
    auction_id          CHAR(26) PRIMARY KEY,
    product_id          CHAR(26)       NOT NULL UNIQUE REFERENCES products (product_id) ON DELETE RESTRICT,

    -- Pricing
    starting_price      DECIMAL(10, 2) NOT NULL CHECK (starting_price > 0),
    reserve_price       DECIMAL(10, 2) CHECK (reserve_price IS NULL OR reserve_price >= starting_price),
    bid_increment       DECIMAL(10, 2) NOT NULL DEFAULT 1.00 CHECK (bid_increment > 0),
    current_price       DECIMAL(10, 2) NOT NULL,

    -- Timing
    start_time          TIMESTAMP      NOT NULL,
    end_time            TIMESTAMP      NOT NULL CHECK (end_time > start_time),
    auto_extend_minutes INTEGER                 DEFAULT 0,           -- Anti-sniping: extend if bid near end

    -- Status
    status              VARCHAR(20)    NOT NULL DEFAULT 'scheduled', -- ('scheduled', 'live', 'ended', 'cancelled')
    -- Winner info (denormalized for performance, set when auction ends)
    winning_bid_id      CHAR(26),
    winning_bidder_id   CHAR(26)       REFERENCES users (user_id) ON DELETE SET NULL,
    final_price         DECIMAL(10, 2),

    created_at          TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL
);

CREATE TABLE bids
(
    bid_id         CHAR(26) PRIMARY KEY,
    auction_id     CHAR(26)       NOT NULL REFERENCES auctions (auction_id) ON DELETE CASCADE,
    bidder_id      CHAR(26)       NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    bid_amount     DECIMAL(10, 2) NOT NULL CHECK (bid_amount > 0),
    max_bid_amount DECIMAL(10, 2),                           -- For proxy bidding
    bid_status     VARCHAR(20)    NOT NULL DEFAULT 'active', -- ('active', 'outbid', 'winning', 'won', 'lost', 'cancelled')
    ip_address     INET,
    user_agent     TEXT,
    created_at     TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- Note: Removed unique constraint - rely on application-level locking
    -- High-frequency inserts benefit from less constraint overhead
    CONSTRAINT bid_amount_positive CHECK (bid_amount > 0)
);

-- Unique index for bid validation (allows better performance than constraint)
CREATE UNIQUE INDEX idx_bids_auction_bidder_time ON bids (auction_id, bidder_id, created_at);

-- ============================================================================
-- SHOPPING CART (UUID for session-based randomness)
-- ============================================================================

CREATE TABLE carts
(
    cart_id    UUID PRIMARY KEY,
    user_id    CHAR(26) UNIQUE REFERENCES users (user_id) ON DELETE CASCADE,
    session_id UUID UNIQUE, -- For guest carts
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    CONSTRAINT user_or_session CHECK (
        (user_id IS NOT NULL AND session_id IS NULL) OR
        (user_id IS NULL AND session_id IS NOT NULL)
        )
);

CREATE TABLE cart_items
(
    cart_item_id   SERIAL PRIMARY KEY,
    cart_id        UUID           NOT NULL REFERENCES carts (cart_id) ON DELETE CASCADE,
    product_id     CHAR(26)       NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    variant_id     INTEGER REFERENCES product_variants (variant_id) ON DELETE CASCADE,
    quantity       INTEGER        NOT NULL CHECK (quantity > 0),
    price_snapshot DECIMAL(10, 2) NOT NULL, -- Price at time of adding
    created_at     TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,

    UNIQUE (cart_id, product_id, variant_id)
);

-- ============================================================================
-- ORDERS & CHECKOUT (ULID for business audit trail)
-- ============================================================================

CREATE TABLE orders
(
    order_id            CHAR(26) PRIMARY KEY,
    user_id             CHAR(26)       NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT,
    order_number        VARCHAR(50)    NOT NULL UNIQUE,                                     -- Human-readable: ORD-2024-123456

    -- Order type
    order_type          VARCHAR(20)    NOT NULL CHECK (order_type IN ('standard', 'auction')),
    auction_id          CHAR(26)       REFERENCES auctions (auction_id) ON DELETE SET NULL, -- If auction order

    -- Pricing
    subtotal            DECIMAL(10, 2) NOT NULL,
    tax_amount          DECIMAL(10, 2) NOT NULL DEFAULT 0,
    shipping_amount     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    discount_amount     DECIMAL(10, 2) NOT NULL DEFAULT 0,
    total_amount        DECIMAL(10, 2) NOT NULL,

    -- Status
    order_status        VARCHAR(20)    NOT NULL DEFAULT 'pending'
        CHECK (order_status IN ('pending', 'processing', 'paid', 'shipped', 'delivered', 'cancelled', 'refunded')),

    -- Addresses (snapshot at order time)
    billing_address_id  INTEGER        REFERENCES user_addresses (address_id) ON DELETE SET NULL,
    shipping_address_id INTEGER        REFERENCES user_addresses (address_id) ON DELETE SET NULL,

    -- Metadata
    customer_notes      TEXT,
    admin_notes         TEXT,
    ip_address          INET,

    created_at          TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL
);

CREATE TABLE order_items
(
    order_item_id SERIAL PRIMARY KEY,
    order_id      CHAR(26)       NOT NULL REFERENCES orders (order_id) ON DELETE CASCADE,
    product_id    CHAR(26)       NOT NULL REFERENCES products (product_id) ON DELETE RESTRICT,
    variant_id    INTEGER        REFERENCES product_variants (variant_id) ON DELETE SET NULL,

    -- Snapshot data (preserve even if product changes/deleted)
    product_name  VARCHAR(255)   NOT NULL,
    sku           VARCHAR(100)   NOT NULL,
    quantity      INTEGER        NOT NULL CHECK (quantity > 0),
    unit_price    DECIMAL(10, 2) NOT NULL,
    subtotal      DECIMAL(10, 2) NOT NULL,
    tax_amount    DECIMAL(10, 2) NOT NULL DEFAULT 0,

    created_at    TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- PAYMENTS (ULID for financial audit trail)
-- ============================================================================

CREATE TABLE payments
(
    payment_id        CHAR(26) PRIMARY KEY,
    order_id          CHAR(26)       NOT NULL REFERENCES orders (order_id) ON DELETE RESTRICT,
    payment_method    VARCHAR(50)    NOT NULL
        CHECK (payment_method IN ('credit_card', 'debit_card', 'paypal', 'stripe', 'bank_transfer', 'other')),

    -- Transaction info
    transaction_id    VARCHAR(255) UNIQUE, -- From payment provider
    amount            DECIMAL(10, 2) NOT NULL,
    currency_code     CHAR(3)        NOT NULL DEFAULT 'USD',

    -- Status
    payment_status    VARCHAR(20)    NOT NULL DEFAULT 'pending'
        CHECK (payment_status IN ('pending', 'processing', 'completed', 'failed', 'refunded', 'cancelled')),

    -- Provider response
    provider_response JSONB,
    failure_reason    TEXT,

    processed_at      TIMESTAMP,
    created_at        TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP NULL
);

-- ============================================================================
-- SHIPPING (SERIAL - internal logistics)
-- ============================================================================

CREATE TABLE shipments
(
    shipment_id        SERIAL PRIMARY KEY,
    order_id           CHAR(26)    NOT NULL REFERENCES orders (order_id) ON DELETE RESTRICT,

    carrier            VARCHAR(100),
    tracking_number    VARCHAR(255),
    shipping_method    VARCHAR(100),
    estimated_delivery DATE,
    actual_delivery    TIMESTAMP,

    shipment_status    VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (shipment_status IN
               ('pending', 'preparing', 'shipped', 'in_transit', 'out_for_delivery', 'delivered', 'failed',
                'returned')),

    shipped_at         TIMESTAMP,
    created_at         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at         TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at         TIMESTAMP NULL
);

CREATE TABLE shipment_tracking_history
(
    tracking_history_id SERIAL PRIMARY KEY,
    shipment_id         INTEGER     NOT NULL REFERENCES shipments (shipment_id) ON DELETE CASCADE,
    status              VARCHAR(20) NOT NULL,
    location            VARCHAR(255),
    description         TEXT,
    tracked_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at          TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- REVIEWS & RATINGS (ULID for user-generated content feed)
-- ============================================================================

CREATE TABLE reviews
(
    review_id            CHAR(26) PRIMARY KEY,
    order_item_id        INTEGER   NOT NULL REFERENCES order_items (order_item_id) ON DELETE CASCADE,
    product_id           CHAR(26)  NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    reviewer_id          CHAR(26)  NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    seller_id            CHAR(26)  NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,

    -- Product review
    product_rating       INTEGER   NOT NULL CHECK (product_rating BETWEEN 1 AND 5),
    product_review_title VARCHAR(255),
    product_review_text  TEXT,

    -- Seller review
    seller_rating        INTEGER   NOT NULL CHECK (seller_rating BETWEEN 1 AND 5),
    seller_review_text   TEXT,

    -- Moderation
    is_verified_purchase BOOLEAN   NOT NULL DEFAULT TRUE,
    is_approved          BOOLEAN   NOT NULL DEFAULT FALSE,
    moderation_notes     TEXT,

    created_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at           TIMESTAMP NULL,

    UNIQUE (order_item_id) -- One review per order item
);

CREATE TABLE review_media
(
    review_media_id SERIAL PRIMARY KEY,
    review_id       CHAR(26)    NOT NULL REFERENCES reviews (review_id) ON DELETE CASCADE,
    media_type      VARCHAR(20) NOT NULL CHECK (media_type IN ('image', 'video')),
    media_url       TEXT        NOT NULL,
    created_at      TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================================
-- WATCHLISTS & WISHLISTS (SERIAL - internal user preferences)
-- ============================================================================

CREATE TABLE watchlists
(
    watchlist_id          SERIAL PRIMARY KEY,
    user_id               CHAR(26)    NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    product_id            CHAR(26)    NOT NULL REFERENCES products (product_id) ON DELETE CASCADE,
    watchlist_type        VARCHAR(20) NOT NULL CHECK (watchlist_type IN ('wishlist', 'auction_watch')),

    -- For auctions: notify preferences
    notify_on_outbid      BOOLEAN              DEFAULT TRUE,
    notify_before_end     BOOLEAN              DEFAULT TRUE,
    notify_minutes_before INTEGER              DEFAULT 30,

    created_at            TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at            TIMESTAMP NULL,

    UNIQUE (user_id, product_id, watchlist_type)
);

-- ============================================================================
-- NOTIFICATIONS (UUID - distributed, non-sequential)
-- ============================================================================

CREATE TABLE notifications
(
    notification_id     UUID PRIMARY KEY,
    user_id             CHAR(26)     NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    notification_type   VARCHAR(50)  NOT NULL, -- ('bid_outbid', 'auction_won', 'auction_ending', 'order_update', 'payment_success', 'payment_failed', 'shipment_update', 'review_reminder', 'system')
    title               VARCHAR(255) NOT NULL,
    message             TEXT         NOT NULL,

    -- Related entities (flexible - can reference any ID type)
    related_entity_type VARCHAR(50),
    related_entity_id   TEXT,                  -- Changed to TEXT to accommodate both ULID and SERIAL
    action_url          TEXT,

    is_read             BOOLEAN      NOT NULL DEFAULT FALSE,
    read_at             TIMESTAMP,

    created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP NULL
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE (including partial indexes for soft delete)
-- ============================================================================

-- Users (ULID - lexicographic ordering)
CREATE INDEX idx_users_email ON users (email);
CREATE INDEX idx_users_username ON users (username);
CREATE INDEX idx_users_account_status ON users (account_status);
CREATE INDEX idx_users_created ON users (user_id DESC);
CREATE INDEX idx_users_deleted_at ON users (deleted_at) WHERE deleted_at IS NULL;

-- User Roles
CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles (role_id);

-- Addresses
CREATE INDEX idx_user_addresses_user_id ON user_addresses (user_id);
CREATE INDEX idx_user_addresses_default ON user_addresses (user_id, is_default) WHERE is_default = TRUE;
CREATE INDEX idx_user_addresses_deleted_at ON user_addresses (deleted_at) WHERE deleted_at IS NULL;

-- Categories
CREATE INDEX idx_categories_parent ON categories (parent_category_id);
CREATE INDEX idx_categories_slug ON categories (slug);
CREATE INDEX idx_categories_deleted_at ON categories (deleted_at) WHERE deleted_at IS NULL;

-- Products (ULID - time-ordered by default)
CREATE INDEX idx_products_seller ON products (seller_id);
CREATE INDEX idx_products_category ON products (category_id);
CREATE INDEX idx_products_type_status ON products (product_type, status);
CREATE INDEX idx_products_sku ON products (sku);
CREATE INDEX idx_products_slug ON products (slug);
CREATE INDEX idx_products_recent ON products (product_id DESC);
CREATE INDEX idx_products_deleted_at ON products (deleted_at) WHERE deleted_at IS NULL;

-- Product Variants
CREATE INDEX idx_product_variants_product ON product_variants (product_id);
CREATE INDEX idx_product_variants_sku ON product_variants (sku);
CREATE INDEX idx_product_variants_deleted_at ON product_variants (deleted_at) WHERE deleted_at IS NULL;

-- Product Media
CREATE INDEX idx_product_media_product ON product_media (product_id, sort_order);

-- Auctions (ULID - time-critical)
CREATE INDEX idx_auctions_product ON auctions (product_id);
CREATE INDEX idx_auctions_status ON auctions (status);
CREATE INDEX idx_auctions_timing ON auctions (start_time, end_time);
CREATE INDEX idx_auctions_ending_soon ON auctions (end_time) WHERE status = 'live';
CREATE INDEX idx_auctions_recent ON auctions (auction_id DESC);
CREATE INDEX idx_auctions_deleted_at ON auctions (deleted_at) WHERE deleted_at IS NULL;

-- Bids (ULID - chronological by nature)
CREATE INDEX idx_bids_auction ON bids (auction_id, bid_id DESC);
CREATE INDEX idx_bids_bidder ON bids (bidder_id);
CREATE INDEX idx_bids_status ON bids (bid_status);
CREATE INDEX idx_bids_recent ON bids (bid_id DESC);

-- Carts (UUID)
CREATE INDEX idx_carts_user ON carts (user_id);
CREATE INDEX idx_carts_session ON carts (session_id);
CREATE INDEX idx_carts_deleted_at ON carts (deleted_at) WHERE deleted_at IS NULL;

-- Cart Items
CREATE INDEX idx_cart_items_cart ON cart_items (cart_id);
CREATE INDEX idx_cart_items_product ON cart_items (product_id);

-- Orders (ULID - business audit trail)
CREATE INDEX idx_orders_user ON orders (user_id);
CREATE INDEX idx_orders_number ON orders (order_number);
CREATE INDEX idx_orders_status ON orders (order_status);
CREATE INDEX idx_orders_recent ON orders (order_id DESC);
CREATE INDEX idx_orders_auction ON orders (auction_id) WHERE auction_id IS NOT NULL;
CREATE INDEX idx_orders_deleted_at ON orders (deleted_at) WHERE deleted_at IS NULL;

-- Order Items
CREATE INDEX idx_order_items_order ON order_items (order_id);
CREATE INDEX idx_order_items_product ON order_items (product_id);

-- Payments (ULID - financial audit)
CREATE INDEX idx_payments_order ON payments (order_id);
CREATE INDEX idx_payments_transaction ON payments (transaction_id);
CREATE INDEX idx_payments_status ON payments (payment_status);
CREATE INDEX idx_payments_recent ON payments (payment_id DESC);
CREATE INDEX idx_payments_deleted_at ON payments (deleted_at) WHERE deleted_at IS NULL;

-- Shipments
CREATE INDEX idx_shipments_order ON shipments (order_id);
CREATE INDEX idx_shipments_tracking ON shipments (tracking_number);
CREATE INDEX idx_shipments_status ON shipments (shipment_status);
CREATE INDEX idx_shipments_deleted_at ON shipments (deleted_at) WHERE deleted_at IS NULL;

-- Reviews (ULID - content feed)
CREATE INDEX idx_reviews_product ON reviews (product_id);
CREATE INDEX idx_reviews_reviewer ON reviews (reviewer_id);
CREATE INDEX idx_reviews_seller ON reviews (seller_id);
CREATE INDEX idx_reviews_approved ON reviews (is_approved);
CREATE INDEX idx_reviews_recent ON reviews (review_id DESC);
CREATE INDEX idx_reviews_deleted_at ON reviews (deleted_at) WHERE deleted_at IS NULL;

-- Watchlists
CREATE INDEX idx_watchlists_user ON watchlists (user_id);
CREATE INDEX idx_watchlists_product ON watchlists (product_id);
CREATE INDEX idx_watchlists_deleted_at ON watchlists (deleted_at) WHERE deleted_at IS NULL;

-- Notifications (UUID - no time ordering needed, use created_at)
CREATE INDEX idx_notifications_user_unread ON notifications (user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_created ON notifications (created_at DESC);
CREATE INDEX idx_notifications_deleted_at ON notifications (deleted_at) WHERE deleted_at IS NULL;

-- Roles and Permissions
CREATE INDEX idx_roles_deleted_at ON roles (deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_permissions_deleted_at ON permissions (deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_tags_deleted_at ON tags (deleted_at) WHERE deleted_at IS NULL;