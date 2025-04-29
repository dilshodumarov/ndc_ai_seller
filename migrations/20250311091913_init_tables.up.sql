CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. client_type
CREATE TABLE IF NOT EXISTS "client_type" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR NOT NULL UNIQUE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. role
CREATE TABLE IF NOT EXISTS "role" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "client_type_id" UUID REFERENCES "client_type"("guid"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3. user
CREATE TABLE IF NOT EXISTS "user" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255),
    "email" VARCHAR(255) UNIQUE,
    "phone_number" VARCHAR(255),
    "password" VARCHAR(255),
    "role_id" UUID REFERENCES "role"("guid") ON DELETE CASCADE,
    "is_active" BOOL DEFAULT FALSE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. business
CREATE TABLE IF NOT EXISTS "business" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "owner_id" UUID REFERENCES "user"("guid") ON DELETE CASCADE, 
    "name" VARCHAR(255),
    "tg_user_name" VARCHAR(255),
    "description" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);

-- 5. integration
CREATE TABLE IF NOT EXISTS "integration" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "owner_id" UUID REFERENCES "business"("guid") ON DELETE CASCADE,
    "integration_token" VARCHAR(255),
    "integration_type" VARCHAR(55),
    "integration_user_name" VARCHAR(255),
    "status" VARCHAR(10) DEFAULT 'active',
    "started_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "stoped_at" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP
);

-- 6. category
CREATE TABLE IF NOT EXISTS "category" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "business_id" UUID REFERENCES "business"("guid"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. product
CREATE TABLE IF NOT EXISTS "product" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "product_id" BIGSERIAL UNIQUE,
    "business_id" UUID REFERENCES "business"("guid"),
    "image_url"   VARCHAR(555),
    "name" VARCHAR(255) NOT NULL,
    "category_id" UUID REFERENCES "category"("guid"),
    "short_info" VARCHAR(255),
    "description" TEXT,
    "cost" INT NOT NULL,
    "count" INT NOT NULL,
    "discount_cost" INT,
    "discount" INT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 9. client
CREATE TABLE IF NOT EXISTS "client" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "platform_id" VARCHAR(255),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255),
    "username" VARCHAR(255) UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "bith_date" VARCHAR,
    "tg_user_name" VARCHAR(255),
    "chat_id" BIGINT,
    "bussnes_id" UUID REFERENCES "business"("guid"),
    "phone" VARCHAR(255),
    "instagram" VARCHAR(255),
    "client_from" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. order
CREATE TABLE IF NOT EXISTS "order" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "order_id" BIGSERIAL UNIQUE,
    "client_id" UUID REFERENCES "client"("guid"),
    "business_id" UUID REFERENCES "business"("guid"),
    "location_url" TEXT,
    "image_url" TEXT,
    "status" VARCHAR(255),
    "total_price" NUMERIC(10, 2) NOT NULL,
    "payment_method" VARCHAR(255) DEFAULT 'online',
    "status_changed_time" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. order_products
CREATE TABLE IF NOT EXISTS "order_products" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "order_id" UUID REFERENCES "order"("guid") ON DELETE CASCADE,
    "product_id" UUID REFERENCES "product"("guid"),
    "count" INT NOT NULL,
    "price" NUMERIC(10, 2) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 12. bot_commands
CREATE TABLE IF NOT EXISTS "bot_commands" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "integration_id" UUID NOT NULL REFERENCES "integration"("guid") ON DELETE CASCADE,
    "command" TEXT NOT NULL, 
    "response" TEXT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 13. chat_history
CREATE TABLE IF NOT EXISTS "chat_history" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "message_id" BIGINT,
    "business_id" UUID REFERENCES "business"("guid"),
    "message" TEXT, 
    "chat_id" BIGINT NOT NULL,
    "platform_id" VARCHAR(255),
    "ai_response" TEXT,
    "reply_to_message_id" BIGINT,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS "telegram_accaunt" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "number" VARCHAR(50) NOT NULL,
    "business_id" UUID REFERENCES "business"("guid"),
    "status" VARCHAR(10) DEFAULT 'active',
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
