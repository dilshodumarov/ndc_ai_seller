CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "client_type" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR NOT NULL UNIQUE,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "role" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "client_type_id" UUID REFERENCES "client_type"("guid"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "user"(
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

CREATE TABLE IF NOT EXISTS  "business" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "owner_id" UUID REFERENCES "user"(guid) ON DELETE CASCADE, 
    "name" VARCHAR(255),
    "integration_token" VARCHAR(255),
    "integration_type" VARCHAR(55),
    "description" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP DEFAULT NULL
);


CREATE TABLE IF NOT EXISTS "client" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "platform_id" VARCHAR(255),
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255),
    "username" VARCHAR(255) UNIQUE,
    "password" VARCHAR(255) NOT NULL,
    "bith_date" VARCHAR,
    "tg_user_name" VARCHAR(255),
    "phone" VARCHAR(255),
    "instagram" VARCHAR(255),
    "client_from" VARCHAR(255),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "category" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL UNIQUE,
    "business_id" UUID REFERENCES "business"("guid"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "attribute" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL,
    "category_id" UUID REFERENCES "category"("guid"),
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "product" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "business_id" UUID REFERENCES "business"("guid"),
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

CREATE TABLE IF NOT EXISTS "integration" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "order" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "client_id" UUID REFERENCES "client"("guid"),
    "integration_id" UUID REFERENCES "integration"("guid"),
    "business_id" UUID REFERENCES "business"("guid"),
    "status" VARCHAR(255),
    "status_changed_time" TIMESTAMP,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "order_products" (
    "guid" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    "order_id" UUID REFERENCES "order"("guid"),
    "product_id" UUID REFERENCES "product"("guid"),
    "count" INT NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
