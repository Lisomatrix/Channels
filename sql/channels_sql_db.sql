--
-- PostgreSQL database dump
--

-- Dumped from database version 13.1
-- Dumped by pg_dump version 13.1

-- Started on 2021-03-15 14:38:37

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- TOC entry 200 (class 1259 OID 49318)
-- Name: App; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."App" (
    "AppID" character varying(150) NOT NULL,
    "Name" character varying(255)
);


--
-- TOC entry 206 (class 1259 OID 57672)
-- Name: Channel; Type: TABLE; Schema: public; Owner: postgres
--

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




--
-- TOC entry 209 (class 1259 OID 57757)
-- Name: Channel_Client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Channel_Client" (
    "clientID" character varying(100) NOT NULL,
    "channelID" bigint NOT NULL
);


--
-- TOC entry 208 (class 1259 OID 57739)
-- Name: Channel_Event; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Channel_Event" (
    "ID" bigint NOT NULL,
    "SenderID" character varying(100) NOT NULL,
    "EventType" character varying(50) NOT NULL,
    "TimeStamp" bigint NOT NULL,
    "Payload" text NOT NULL,
    "ChannelID" bigint NOT NULL
);


--
-- TOC entry 207 (class 1259 OID 57737)
-- Name: Channel_Event_ID_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Channel_Event_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3052 (class 0 OID 0)
-- Dependencies: 207
-- Name: Channel_Event_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Channel_Event_ID_seq" OWNED BY public."Channel_Event"."ID";


--
-- TOC entry 205 (class 1259 OID 57670)
-- Name: Channel_ID_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."Channel_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- TOC entry 3053 (class 0 OID 0)
-- Dependencies: 205
-- Name: Channel_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."Channel_ID_seq" OWNED BY public."Channel"."ID";


--
-- TOC entry 201 (class 1259 OID 49323)
-- Name: Client; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Client" (
    "ID" character varying(100) NOT NULL,
    "Username" character varying(100),
    "AppID" character varying(155) NOT NULL,
    "Extra" text
);

--
-- TOC entry 202 (class 1259 OID 49430)
-- Name: Device; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public."Device" (
    "ID" character varying(50) NOT NULL,
    "Token" character varying(100) NOT NULL,
    "ClientID" character varying(100) NOT NULL
);


--
-- TOC entry 204 (class 1259 OID 57514)
-- Name: NewChannel; Type: TABLE; Schema: public; Owner: postgres
--

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


--
-- TOC entry 203 (class 1259 OID 57512)
-- Name: NewChannel_ID_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public."NewChannel_ID_seq"
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;




--
-- TOC entry 3054 (class 0 OID 0)
-- Dependencies: 203
-- Name: NewChannel_ID_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public."NewChannel_ID_seq" OWNED BY public."NewChannel"."ID";


--
-- TOC entry 2888 (class 2604 OID 57675)
-- Name: Channel ID; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel" ALTER COLUMN "ID" SET DEFAULT nextval('public."Channel_ID_seq"'::regclass);


--
-- TOC entry 2889 (class 2604 OID 57742)
-- Name: Channel_Event ID; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Event" ALTER COLUMN "ID" SET DEFAULT nextval('public."Channel_Event_ID_seq"'::regclass);


--
-- TOC entry 2882 (class 2604 OID 57517)
-- Name: NewChannel ID; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."NewChannel" ALTER COLUMN "ID" SET DEFAULT nextval('public."NewChannel_ID_seq"'::regclass);


--
-- TOC entry 2891 (class 2606 OID 49322)
-- Name: App App_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."App"
    ADD CONSTRAINT "App_pkey" PRIMARY KEY ("AppID");


--
-- TOC entry 2906 (class 2606 OID 57744)
-- Name: Channel_Event Channel_Event_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Event"
    ADD CONSTRAINT "Channel_Event_pkey" PRIMARY KEY ("ID");


--
-- TOC entry 2901 (class 2606 OID 57680)
-- Name: Channel Channel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT "Channel_pkey" PRIMARY KEY ("ID");


--
-- TOC entry 2893 (class 2606 OID 49330)
-- Name: Client Client_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Client"
    ADD CONSTRAINT "Client_pkey" PRIMARY KEY ("ID");


--
-- TOC entry 2895 (class 2606 OID 49434)
-- Name: Device Device_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Device"
    ADD CONSTRAINT "Device_pkey" PRIMARY KEY ("ID");


--
-- TOC entry 2897 (class 2606 OID 57522)
-- Name: NewChannel NewChannel_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT "NewChannel_pkey" PRIMARY KEY ("ID");


--
-- TOC entry 2899 (class 2606 OID 57610)
-- Name: NewChannel app_ch_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT app_ch_unique UNIQUE ("ChannelID", "AppID");


--
-- TOC entry 2909 (class 2606 OID 57771)
-- Name: Channel_Client client_channelID_unique; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT "client_channelID_unique" UNIQUE ("clientID", "channelID");


--
-- TOC entry 2904 (class 2606 OID 57682)
-- Name: Channel unique_app_channel; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT unique_app_channel UNIQUE ("AppID", "ChannelID");


--
-- TOC entry 2902 (class 1259 OID 57777)
-- Name: appID_channelID_indexx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "appID_channelID_indexx" ON public."Channel" USING btree ("ChannelID", "AppID");


--
-- TOC entry 2907 (class 1259 OID 57779)
-- Name: channelID_TimeStamp_Indexx; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX "channelID_TimeStamp_Indexx" ON public."Channel_Event" USING btree ("ChannelID", "TimeStamp");


--
-- TOC entry 2912 (class 2606 OID 57619)
-- Name: NewChannel ch_appID_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."NewChannel"
    ADD CONSTRAINT "ch_appID_fk" FOREIGN KEY ("AppID") REFERENCES public."App"("AppID") ON UPDATE CASCADE ON DELETE CASCADE NOT VALID;


--
-- TOC entry 2914 (class 2606 OID 57772)
-- Name: Channel_Event channelID_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Event"
    ADD CONSTRAINT "channelID_fk" FOREIGN KEY ("ChannelID") REFERENCES public."Channel"("ID") NOT VALID;


--
-- TOC entry 2916 (class 2606 OID 57765)
-- Name: Channel_Client channel_client_channel_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT channel_client_channel_fk FOREIGN KEY ("channelID") REFERENCES public."Channel"("ID");


--
-- TOC entry 2915 (class 2606 OID 57760)
-- Name: Channel_Client client_channel_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel_Client"
    ADD CONSTRAINT client_channel_fk FOREIGN KEY ("clientID") REFERENCES public."Client"("ID");


--
-- TOC entry 2913 (class 2606 OID 57683)
-- Name: Channel fk_channel_app; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Channel"
    ADD CONSTRAINT fk_channel_app FOREIGN KEY ("AppID") REFERENCES public."App"("AppID");


--
-- TOC entry 2910 (class 2606 OID 49331)
-- Name: Client fk_client_app; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Client"
    ADD CONSTRAINT fk_client_app FOREIGN KEY ("AppID") REFERENCES public."App"("AppID");


--
-- TOC entry 2911 (class 2606 OID 49440)
-- Name: Device fk_device_client; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public."Device"
    ADD CONSTRAINT fk_device_client FOREIGN KEY ("ClientID") REFERENCES public."Client"("ID") NOT VALID;


-- Completed on 2021-03-15 14:38:38

--
-- PostgreSQL database dump complete
--

