User

- id (PK, uuid)
- email (unique, not null)
- firstName (not null)
- lastName (nullable)
- emailVerified (timestamp, nullable)
- createdAt (timestamp)
- updatedAt (timestamp)

RELATIONS:

- has many → Account
- has many → Membership

---

Account

- id (PK, uuid)
- userId (FK → User.id, not null)
- provider (enum: email, google, github, microsoft)
  - providerAccountId (not null)
- passwordHash (nullable)
- accessToken (nullable)
- refreshToken (nullable)
- tokenExpiresAt (timestamp, nullable)
- createdAt (timestamp)
- updatedAt (timestamp)

CONSTRAINTS:

- unique(provider, providerAccountId)

RELATIONS:

- belongs to → User

---

Organization

- id (PK, uuid)
- name (not null)
- slug (unique, not null)
- createdAt (timestamp)
- updatedAt (timestamp)

RELATIONS:

- has many → Membership
- has many → ApiKey
- has many → Message

---

Membership

- id (PK, uuid)
- userId (FK → User.id, not null)
- orgId (FK → Organization.id, not null)
- role (enum: admin, member)
- createdAt (timestamp)
- updatedAt (timestamp)

CONSTRAINTS:

- unique(userId, orgId)

RELATIONS:

- belongs to → User
- belongs to → Organization

---

ApiKey

- id (PK, uuid)
- orgId (FK → Organization.id, not null)
- name (not null)
- hashedKey (unique, not null)
- prefix (not null)
- expiresAt (timestamp, nullable)
- lastUsedAt (timestamp, nullable)
- createdAt (timestamp)
- updatedAt (timestamp)

RELATIONS:

- belongs to → Organization

---

Message

- id (PK, uuid)
- orgId (FK → Organization.id, not null)
- url (not null)
- method (not null, default: POST)
- payload (jsonb, not null)
- status (enum: pending, success, failed)
- createdAt (timestamp)
- updatedAt (timestamp)

RELATIONS:

- belongs to → Organization
- has many → DeliveryAttempt

---

DeliveryAttempt

- id (PK, uuid)
- messageId (FK → Message.id, not null)
- attemptNumber (int, not null)
- statusCode (int, nullable)
- responseBody (text, nullable)
- errorMessage (text, nullable)
- durationMs (int, nullable)
- attemptedAt (timestamp, not null)

RELATIONS:

- belongs to → Message
