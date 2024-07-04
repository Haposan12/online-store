CREATE TABLE "public"."product" (
  "id" serial8,
  "name" varchar(50),
  "description" varchar(100),
  "category_id" int8,
  "price" float8,
  "stock" int2,
  "created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_category" FOREIGN KEY ("category_id") REFERENCES "public"."category" ("id")
);


CREATE TABLE "public"."category" (
  "id" serial8,
  "name" varchar(50),
  "description" varchar(100),
  "created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("id")
);

CREATE TABLE "public"."customer" (
  "customer_id" serial8,
  "first_name" varchar(50),
  "last_name" varchar(50),
  "email" varchar(50),
  "password" varchar(50),
  "address" varchar(50),
  "phone_number" varchar(50),
  "created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("customer_id")
);


CREATE TABLE "public"."cart" (
 "cart_id" serial8,
 "quantity" int8,
 "customer_id" int8,
 "product_id" int8,
"created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("cart_id"),
  CONSTRAINT "fk_customer" FOREIGN KEY ("customer_id") REFERENCES "public"."customer" ("customer_id"),
  CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "public"."product" ("id")
);

CREATE TABLE "public"."order" (
 "id" serial8,
 "total_price" float8,
 "customer_id" int8,
 "payment_id" int8 default NULL,
"created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_customer" FOREIGN KEY ("customer_id") REFERENCES "public"."customer" ("customer_id"),
  CONSTRAINT "fk_payment" FOREIGN KEY ("payment_id") REFERENCES "public"."payment" ("id")
);

CREATE TABLE "public"."order_item" (
 "id" serial8,
 "quantity" int8,
 "product_id" int8,
 "price" float8,
 "order_id" int8 default NULL,
"created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order" FOREIGN KEY ("order_id") REFERENCES "public"."order" ("id"),
  CONSTRAINT "fk_product" FOREIGN KEY ("product_id") REFERENCES "public"."product" ("id")
);

CREATE TABLE "public"."payment" (
 "id" serial8,
 "method" varchar(100),
 "amount" float8,
 "status" varchar(50),
"created_at" timestamptz(6) DEFAULT now(),
  "created_by" varchar(50) DEFAULT 'system',
  "updated_at" timestamptz(6),
  "updated_by" varchar(50),
  "deleted_at" timestamptz(6),
  "deleted_by" varchar(50),
  PRIMARY KEY ("id")
);