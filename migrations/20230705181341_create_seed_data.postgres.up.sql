INSERT INTO public.agent (username,password,name,email,phone,is_owner,created_at,updated_at) VALUES
	 ('test','test','test','test@test.com','test',true,'2023-05-07 00:00:00.000','2023-05-07 00:00:00.000'),
	 ('test2','test2','test2','test2@test2.com','test2',false,'2023-05-07 00:00:00.000','2023-05-07 00:00:00.000');
INSERT INTO public.balance (bal_cash,bal_qr,created_at,updated_at) VALUES
	 ('sb=500&eb=1270','sb=100&eb=561','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000'),
	 ('sb=1159&eb=1791','sb=561&eb=241','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000');
INSERT INTO public.item (name,des,created_at,updated_at) VALUES
	 ('MD2 Juice','Freshly pressed high quality pineapple juice, 100% natural with no sugar added','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000'),
	 ('MD2 Sliced Fruit','Freshly cut fruit, you can taste its original sweetness','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000');
INSERT INTO public.inventory (start_item_bal,end_item_bal,created_at,updated_at) VALUES
	 ('0=150','0=5','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000'),
	 ('0=520&1=50','0=14&1=15','2023-06-07 00:00:00.000','2023-06-07 00:00:00.000');
INSERT INTO public.operation (start_time,end_time,location,agent_id,total_sales_qty,total_cost,total_sales_amount,net_profit,balance_id,inventory_id,created_at,updated_at) VALUES
	 ('2023-01-07 00:00:00','2023-01-07 00:00:00','Test Location',1,150,330.89,1200,869.11,1,1,'2023-06-07 00:00:00','2023-06-07 00:00:00');
INSERT INTO public.sale (amount,quantity,payment_type,operation_id,item_id,created_at,updated_at) VALUES
	 (32,4,1,1,1,'2023-06-07 00:00:00.000','2023-06-07 00:00:00.000'),
	 (10,1,1,1,2,'2023-06-07 00:00:00.000','2023-06-07 00:00:00.000');