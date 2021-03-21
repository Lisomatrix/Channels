-- POSTGRESQL

CREATE TABLE public."App" (
    "AppID" character varying(150) NOT NULL,
    "Name" character varying(255)
);

CREATE TABLE public."Channel" (
    "ID" bigint NOT NULL,
    "ChannelID" character varying(100) NOT NULL,
    "AppID" character varying(150),
    "Name" character varying(150),
    "Created_At" bigint,
    "IsClosed" boolean,
    "Extra" text,
    "Persistent" boolean,
    "Private" boolean,
    "Presence" boolean
);

CREATE TABLE public."Channel_Client" (
    "clientID" character varying(100) NOT NULL,
    "channelID" bigint NOT NULL
);


CREATE TABLE public."Channel_Event" (
    "ID" bigint NOT NULL,
    "SenderID" character varying(100) NOT NULL,
    "EventType" character varying(50) NOT NULL,
    "TimeStamp" bigint NOT NULL,
    "Payload" text NOT NULL,
    "ChannelID" bigint NOT NULL
);

CREATE SEQUENCE public."Channel_Event_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Channel_Event_ID_seq" OWNED BY public."Channel_Event"."ID";

CREATE SEQUENCE public."Channel_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public."Channel_ID_seq" OWNED BY public."Channel"."ID";

CREATE TABLE public."Client" (
    "ID" character varying(100) NOT NULL,
    "Username" character varying(100),
    "AppID" character varying(155) NOT NULL,
    "Extra" text
);

CREATE TABLE public."Device" (
    "ID" character varying(50) NOT NULL,
    "Token" character varying(100) NOT NULL,
    "ClientID" character varying(100) NOT NULL
);


CREATE TABLE public."NewChannel" (
    "ID" bigint NOT NULL,
    "ChannelID" character varying(150) NOT NULL,
    "AppID" character varying(150) NOT NULL,
    "Created_At" bigint DEFAULT 0 NOT NULL,
    "IsClosed" boolean DEFAULT false NOT NULL,
    "Extra" text NOT NULL,
    "Persistent" boolean DEFAULT false NOT NULL,
    "Private" boolean DEFAULT false NOT NULL,
    "Presence" boolean DEFAULT false NOT NULL
);

CREATE SEQUENCE public."NewChannel_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE public."NewChannel_ID_seq" OWNED BY public."NewChannel"."ID";

ALTER TABLE ONLY public."Channel" ALTER COLUMN "ID" SET DEFAULT nextval('public."Channel_ID_seq"'::regclass);

ALTER TABLE ONLY public."Channel_Event" ALTER COLUMN "ID" SET DEFAULT nextval('public."Channel_Event_ID_seq"'::regclass);

ALTER TABLE ONLY public."NewChannel" ALTER COLUMN "ID" SET DEFAULT nextval('public."NewChannel_ID_seq"'::regclass);

ALTER TABLE ONLY public."App"
    ADD CONSTRAINT "App_pkey" PRIMARY KEY ("AppID");


ALTER TABLE ONLY public."Channel_Event"
    ADD CONSTRAINT "Channel_Event_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT "Channel_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY public."Client"
    ADD CONSTRAINT "Client_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY public."Device"
    ADD CONSTRAINT "Device_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT "NewChannel_pkey" PRIMARY KEY ("ID");

ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT app_ch_unique UNIQUE ("ChannelID", "AppID");

ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT "client_channelID_unique" UNIQUE ("clientID", "channelID");

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT unique_app_channel UNIQUE ("AppID", "ChannelID");


CREATE INDEX "appID_channelID_indexx" ON public."Channel" USING btree ("ChannelID", "AppID");


CREATE INDEX "channelID_TimeStamp_Indexx" ON public."Channel_Event" USING btree ("ChannelID", "TimeStamp");


ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT "ch_appID_fk" FOREIGN KEY ("AppID") REFERENCES public."App"("AppID") ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;

ALTER TABLE ONLY public."Channel_Event"
    ADD CONSTRAINT "channelID_fk" FOREIGN KEY ("ChannelID") REFERENCES public."Channel"("ID") NOT VALID;

ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT channel_client_channel_fk FOREIGN KEY ("channelID") REFERENCES public."Channel"("ID");


ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT client_channel_fk FOREIGN KEY ("clientID") REFERENCES public."Client"("ID");

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT fk_channel_app FOREIGN KEY ("AppID") REFERENCES public."App"("AppID");

ALTER TABLE ONLY public."Client"
    ADD CONSTRAINT fk_client_app FOREIGN KEY ("AppID") REFERENCES public."App"("AppID");

ALTER TABLE ONLY public."Device"
    ADD CONSTRAINT fk_device_client FOREIGN KEY ("ClientID") REFERENCES public."Client"("ID") NOT VALID;
