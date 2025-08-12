ALTER TABLE badge_description
  ADD COLUMN description TEXT;

ALTER TABLE public.badge_description_version
  ADD COLUMN description TEXT;
