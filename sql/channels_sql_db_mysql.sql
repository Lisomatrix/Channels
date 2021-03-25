CREATE TABLE App (
    AppID character varying(150) NOT NULL,
    Name character varying(255)
);

CREATE TABLE Channel (
    ID bigint NOT NULL AUTO_INCREMENT,
    ChannelID character varying(100) NOT NULL,
    AppID character varying(150),
    Name character varying(150),
    Created_At bigint,
    IsClosed boolean,
    Extra text,
    Persistent boolean,
    Private boolean,
    Presence boolean,
    Publish boolean,
    primary key (ID)
);

CREATE TABLE Channel_Client (
    clientID character varying(100) NOT NULL,
    channelID bigint NOT NULL
);

CREATE TABLE Channel_Event (
    ID bigint NOT NULL AUTO_INCREMENT,
    SenderID character varying(100) NOT NULL,
    EventType character varying(50) NOT NULL,
    TimeStamp bigint NOT NULL,
    Payload text NOT NULL,
    ChannelID bigint NOT NULL,
    primary key (ID)
);

CREATE TABLE Client (
    ID character varying(100) NOT NULL,
    Username character varying(100),
    AppID character varying(155) NOT NULL,
    Extra text
);

CREATE TABLE Device (
    ID character varying(50) NOT NULL,
    Token character varying(350) NOT NULL,
    ClientID character varying(100) NOT NULL
);

ALTER TABLE App ADD CONSTRAINT App_pkey PRIMARY KEY (AppID);

ALTER TABLE Client ADD CONSTRAINT Client_pkey PRIMARY KEY (ID);

ALTER TABLE Device ADD CONSTRAINT Device_pkey PRIMARY KEY (ID);

ALTER TABLE Channel_Client ADD CONSTRAINT client_channelID_unique UNIQUE (clientID, channelID);

ALTER TABLE Channel ADD CONSTRAINT unique_app_channel UNIQUE (AppID, ChannelID);

CREATE INDEX appID_channelID_indexx ON Channel (ChannelID, AppID);

CREATE INDEX channelID_TimeStamp_Indexx ON Channel_Event (ChannelID, TimeStamp);

ALTER TABLE Channel_Event ADD CONSTRAINT channelID_fk FOREIGN KEY (ChannelID) REFERENCES Channel(ID);

ALTER TABLE Channel_Client ADD CONSTRAINT channel_client_channel_fk FOREIGN KEY (channelID) REFERENCES Channel(ID);

ALTER TABLE Channel_Client ADD CONSTRAINT client_channel_fk FOREIGN KEY (clientID) REFERENCES Client(ID);

ALTER TABLE Channel ADD CONSTRAINT fk_channel_app FOREIGN KEY (AppID) REFERENCES App(AppID);

ALTER TABLE Client ADD CONSTRAINT fk_client_app FOREIGN KEY (AppID) REFERENCES App(AppID);

ALTER TABLE Device ADD CONSTRAINT fk_device_client FOREIGN KEY (ClientID) REFERENCES Client(ID);