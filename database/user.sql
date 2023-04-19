drop user if exists app;

create user app password '0u1Ot1m3';

GRANT ALL ON SCHEMA outatime TO app;

GRANT EXECUTE ON FUNCTION outatime.make_user_page_func() TO app;

GRANT ALL ON SEQUENCE outatime.event_event_id_seq TO app;

GRANT ALL ON SEQUENCE outatime.link_link_id_seq TO app;

GRANT ALL ON SEQUENCE outatime.subscriber_sub_id_seq TO app;

GRANT ALL ON SEQUENCE outatime.subscription_subscription_id_seq TO app;

GRANT ALL ON SEQUENCE outatime.user_page_page_id_seq TO app;

GRANT ALL ON SEQUENCE outatime.user_user_id_seq TO app;

GRANT ALL ON TABLE outatime.event TO app;

GRANT ALL ON TABLE outatime.link TO app;

GRANT ALL ON TABLE outatime.subscriber TO app;

GRANT ALL ON TABLE outatime.subscription TO app;

GRANT ALL ON TABLE outatime."user" TO app;

GRANT ALL ON TABLE outatime.user_page TO app;

