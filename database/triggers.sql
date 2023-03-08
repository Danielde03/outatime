-- Create a user-page entry when a new user is made
CREATE OR REPLACE FUNCTION outatime.make_user_page_func() RETURNS trigger
   LANGUAGE plpgsql AS
$$BEGIN
   
    INSERT INTO outatime.user_page (user_id, banner)
        VALUES (NEW.user_id, ' ');

    RETURN NEW;

END;$$;

CREATE OR REPLACE TRIGGER make_user_page_air
AFTER INSERT ON "outatime"."user"
FOR EACH ROW
EXECUTE FUNCTION outatime.make_user_page_func();
