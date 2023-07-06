--
-- PostgreSQL database dump
--

-- Dumped from database version 15.3 (Ubuntu 15.3-1.pgdg22.04+1)
-- Dumped by pg_dump version 15.3 (Ubuntu 15.3-1.pgdg22.04+1)

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
-- Name: agent; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.agent (
    id integer NOT NULL,
    username character varying(255) DEFAULT ''::character varying NOT NULL,
    password character varying(255) DEFAULT ''::character varying NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    phone character varying(255) DEFAULT ''::character varying NOT NULL,
    is_owner boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.agent OWNER TO sp3ck;

--
-- Name: agent_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.agent_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.agent_id_seq OWNER TO sp3ck;

--
-- Name: agent_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.agent_id_seq OWNED BY public.agent.id;


--
-- Name: balance; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.balance (
    id integer NOT NULL,
    bal_cash character varying(255) DEFAULT 'sb=0&eb=0'::character varying NOT NULL,
    bal_qr character varying(255) DEFAULT 'sb=0&eb=0'::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.balance OWNER TO sp3ck;

--
-- Name: balance_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.balance_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.balance_id_seq OWNER TO sp3ck;

--
-- Name: balance_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.balance_id_seq OWNED BY public.balance.id;


--
-- Name: inventory; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.inventory (
    id integer NOT NULL,
    start_item_bal character varying(255) DEFAULT '1=0'::character varying NOT NULL,
    end_item_bal character varying(255) DEFAULT '1=0'::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.inventory OWNER TO sp3ck;

--
-- Name: inventory_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.inventory_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.inventory_id_seq OWNER TO sp3ck;

--
-- Name: inventory_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.inventory_id_seq OWNED BY public.inventory.id;


--
-- Name: item; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.item (
    id integer NOT NULL,
    name character varying(255) DEFAULT ''::character varying NOT NULL,
    des character varying(255) DEFAULT ''::character varying NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.item OWNER TO sp3ck;

--
-- Name: item_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.item_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.item_id_seq OWNER TO sp3ck;

--
-- Name: item_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.item_id_seq OWNED BY public.item.id;


--
-- Name: operation; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.operation (
    id integer NOT NULL,
    start_time timestamp without time zone NOT NULL,
    end_time timestamp without time zone NOT NULL,
    location character varying(255) DEFAULT ''::character varying NOT NULL,
    agent_id integer NOT NULL,
    total_sales_qty integer DEFAULT 0 NOT NULL,
    total_cost numeric DEFAULT 0.00 NOT NULL,
    total_sales_amount numeric DEFAULT 0.00 NOT NULL,
    net_profit numeric DEFAULT 0.00 NOT NULL,
    balance_id integer NOT NULL,
    inventory_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.operation OWNER TO sp3ck;

--
-- Name: operation_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.operation_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.operation_id_seq OWNER TO sp3ck;

--
-- Name: operation_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.operation_id_seq OWNED BY public.operation.id;


--
-- Name: sale; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.sale (
    id integer NOT NULL,
    amount numeric DEFAULT 0.0 NOT NULL,
    quantity integer DEFAULT 0 NOT NULL,
    payment_type integer DEFAULT 1 NOT NULL,
    operation_id integer NOT NULL,
    item_id integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


ALTER TABLE public.sale OWNER TO sp3ck;

--
-- Name: sale_id_seq; Type: SEQUENCE; Schema: public; Owner: sp3ck
--

CREATE SEQUENCE public.sale_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.sale_id_seq OWNER TO sp3ck;

--
-- Name: sale_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: sp3ck
--

ALTER SEQUENCE public.sale_id_seq OWNED BY public.sale.id;


--
-- Name: schema_migration; Type: TABLE; Schema: public; Owner: sp3ck
--

CREATE TABLE public.schema_migration (
    version character varying(14) NOT NULL
);


ALTER TABLE public.schema_migration OWNER TO sp3ck;

--
-- Name: agent id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.agent ALTER COLUMN id SET DEFAULT nextval('public.agent_id_seq'::regclass);


--
-- Name: balance id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.balance ALTER COLUMN id SET DEFAULT nextval('public.balance_id_seq'::regclass);


--
-- Name: inventory id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.inventory ALTER COLUMN id SET DEFAULT nextval('public.inventory_id_seq'::regclass);


--
-- Name: item id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.item ALTER COLUMN id SET DEFAULT nextval('public.item_id_seq'::regclass);


--
-- Name: operation id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.operation ALTER COLUMN id SET DEFAULT nextval('public.operation_id_seq'::regclass);


--
-- Name: sale id; Type: DEFAULT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.sale ALTER COLUMN id SET DEFAULT nextval('public.sale_id_seq'::regclass);


--
-- Name: agent agent_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.agent
    ADD CONSTRAINT agent_pkey PRIMARY KEY (id);


--
-- Name: balance balance_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.balance
    ADD CONSTRAINT balance_pkey PRIMARY KEY (id);


--
-- Name: inventory inventory_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.inventory
    ADD CONSTRAINT inventory_pkey PRIMARY KEY (id);


--
-- Name: item item_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.item
    ADD CONSTRAINT item_pkey PRIMARY KEY (id);


--
-- Name: operation operation_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.operation
    ADD CONSTRAINT operation_pkey PRIMARY KEY (id);


--
-- Name: sale sale_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.sale
    ADD CONSTRAINT sale_pkey PRIMARY KEY (id);


--
-- Name: schema_migration schema_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.schema_migration
    ADD CONSTRAINT schema_migration_pkey PRIMARY KEY (version);


--
-- Name: schema_migration_version_idx; Type: INDEX; Schema: public; Owner: sp3ck
--

CREATE UNIQUE INDEX schema_migration_version_idx ON public.schema_migration USING btree (version);


--
-- Name: operation operation_agent_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.operation
    ADD CONSTRAINT operation_agent_id_fk FOREIGN KEY (agent_id) REFERENCES public.agent(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: operation operation_balance_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.operation
    ADD CONSTRAINT operation_balance_id_fk FOREIGN KEY (balance_id) REFERENCES public.balance(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: operation operation_inventory_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.operation
    ADD CONSTRAINT operation_inventory_id_fk FOREIGN KEY (inventory_id) REFERENCES public.inventory(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: sale sale_item_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.sale
    ADD CONSTRAINT sale_item_id_fk FOREIGN KEY (item_id) REFERENCES public.item(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- Name: sale sale_operation_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: sp3ck
--

ALTER TABLE ONLY public.sale
    ADD CONSTRAINT sale_operation_id_fk FOREIGN KEY (operation_id) REFERENCES public.operation(id) ON UPDATE CASCADE ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

