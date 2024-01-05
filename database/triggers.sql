-- Create a user-page entry when a new user is made
CREATE OR REPLACE FUNCTION outatime.make_user_page_func() RETURNS trigger
   LANGUAGE plpgsql AS
$$BEGIN
   
    INSERT INTO outatime.user_page (user_id, "aboutUs", banner)
        VALUES (NEW.user_id, ' ', ' ');

    RETURN NEW;

END;$$;

CREATE OR REPLACE TRIGGER make_user_page_air
AFTER INSERT ON "outatime"."user"
FOR EACH ROW
EXECUTE FUNCTION outatime.make_user_page_func();



-- Update event number when new events are created
CREATE OR REPLACE FUNCTION outatime.add_to_user_event_count() RETURNS trigger
   LANGUAGE plpgsql AS
$$BEGIN
   
    UPDATE "outatime"."user" SET events = (
        SELECT events FROM "outatime"."user" WHERE user_id = NEW.user_id
    ) + 1
     WHERE user_id = NEW.user_id;

    RETURN NEW;

END;$$;

CREATE OR REPLACE TRIGGER add_user_count_air
AFTER INSERT ON "outatime"."event"
FOR EACH ROW
EXECUTE FUNCTION outatime.add_to_user_event_count();


-- Update event number when events are removed
CREATE OR REPLACE FUNCTION outatime.remove_from_user_event_count() RETURNS trigger
   LANGUAGE plpgsql AS
$$BEGIN
   
    UPDATE "outatime"."user" SET events = (
        SELECT events FROM "outatime"."user" WHERE user_id = NEW.user_id
    ) - 1
     WHERE user_id = NEW.user_id;

    RETURN NEW;

END;$$;

CREATE OR REPLACE TRIGGER remove_user_count_air
AFTER DELETE ON "outatime"."event"
FOR EACH ROW
EXECUTE FUNCTION outatime.remove_from_user_event_count();