# Referral System Implementation

## Overview

The referral system allows users to share links that, when followed by others, can track donations back to the original sharer. This creates a social incentive for users to promote tree planting and track their impact through referrals.

## How It Works

### 1. Creating Share Links

Users can create two types of share links:

- **Profile Share**: A general link to their profile (`POST /v1/shares/profile`)
- **Donation Share**: A specific link to a particular donation (`POST /v1/shares/donation`)

Each share link generates a unique slug that can be shared publicly.

### 2. Following Referral Links

When someone visits a share link (e.g., `GET /v1/shares/resolve/{slug}`), they receive information about:
- The type of share (profile or donation)
- The referral user ID (who created the share)
- Any associated reference data

### 3. Making Referred Donations

When making a donation through a referral link, the frontend should:
1. Extract the `referral_user_id` from the share resolution
2. Include it in the payment intent request (`POST /v1/payments/intents`)
3. The referral user ID gets stored in the payment metadata

### 4. Tracking Referrals

The system automatically:
- Stores the `referral_user_id` in the donation record
- Maintains referential integrity with foreign key constraints
- Provides referral statistics for users

## API Endpoints

### Share Management
- `POST /v1/shares/profile` - Create profile share link
- `POST /v1/shares/donation` - Create donation share link
- `GET /v1/shares/resolve/{slug}` - Resolve share link (public)
- `GET /v1/shares` - List user's share tokens
- `DELETE /v1/shares/{id}` - Delete share token
- `GET /v1/shares/stats` - Get referral statistics

### Payment with Referral
- `POST /v1/payments/intents` - Create payment intent with optional `referral_user_id`

## Database Schema

### Donations Table
```sql
ALTER TABLE donations ADD COLUMN referral_user_id text;
CREATE INDEX idx_donations_referral_user ON donations(referral_user_id);
ALTER TABLE donations ADD CONSTRAINT fk_donations_referral_user 
    FOREIGN KEY (referral_user_id) REFERENCES user_auth(auth_user_id) ON DELETE SET NULL;
```

### Share Tokens Table
```sql
CREATE TABLE share_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    auth_user_id text NOT NULL,
    kind text NOT NULL, -- 'profile' or 'donation'
    ref_id uuid, -- NULL for profile, donation ID for donation shares
    slug text UNIQUE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
```

## Example Flow

1. **User A creates a profile share:**
   ```bash
   POST /v1/shares/profile
   # Returns: {"slug": "john-doe-profile", "url": "http://localhost:8080/share/john-doe-profile"}
   ```

2. **User B visits the share link:**
   ```bash
   GET /v1/shares/resolve/john-doe-profile
   # Returns: {"referral_user_id": "user-a-id", "kind": "profile"}
   ```

3. **User B makes a donation with referral:**
   ```bash
   POST /v1/payments/intents
   {
     "provider": "cloudpayments",
     "amount_minor": 1900,
     "currency": "RUB",
     "referral_user_id": "user-a-id",
     "success_return_url": "https://example.com/success",
     "fail_return_url": "https://example.com/fail"
   }
   ```

4. **System creates donation with referral tracking:**
   - Payment created with `referral_user_id` in metadata
   - Donation created with `referral_user_id` field populated
   - User A can see this referral in their statistics

## Referral Statistics

Users can view their referral impact:
```bash
GET /v1/shares/stats
# Returns:
{
  "total_referrals": 5,
  "total_trees_planted": 25,
  "recent_referrals": [...]
}
```

## Security Considerations

- Share tokens are public but don't expose sensitive user information
- Referral tracking is optional and transparent to donors
- Users can delete their share tokens at any time
- Foreign key constraints ensure data integrity

## Frontend Integration

The frontend should:
1. Store referral information from share links in session/local storage
2. Include referral_user_id in payment forms when available
3. Display referral statistics to users
4. Provide easy sharing functionality for profiles and donations

## Benefits

- **Social Proof**: Users can see the impact of their sharing
- **Community Building**: Encourages users to promote the platform
- **Transparency**: Clear tracking of referral relationships
- **Gamification**: Users can compete on referral impact
- **Analytics**: Platform can measure viral growth and user engagement
