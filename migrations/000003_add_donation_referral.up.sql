-- Add referral_user_id field to donations table
-- This tracks which user referred the donation through a shared link

ALTER TABLE donations ADD COLUMN referral_user_id text;
CREATE INDEX idx_donations_referral_user ON donations(referral_user_id);

-- Add foreign key constraint to ensure referral_user_id references a valid user
ALTER TABLE donations ADD CONSTRAINT fk_donations_referral_user 
    FOREIGN KEY (referral_user_id) REFERENCES user_auth(auth_user_id) ON DELETE SET NULL;

-- Ensure share_tokens slug uniqueness at database level (additional safety)
-- This constraint already exists from the initial migration, but let's verify it's enforced
ALTER TABLE share_tokens ADD CONSTRAINT IF NOT EXISTS share_tokens_slug_unique UNIQUE (slug);
