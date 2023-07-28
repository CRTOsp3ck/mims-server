##### MIMS API SERVER (DATASTORE)  ##### 
Food for thought: Maybe I should just make this a datastore microservice and create a seperate api server microservice when shit starts to hit the fan.

##### README #####
1. Install go
2. Install postgres
3. add root as psql user
4. set root psql password
5. create db with name root
6. git clone mims-server
7. nano main.go db settings
8. nano database.yml settings
9. login to psql
10. add seed data
11. i know it is shit systems design atm, but its a start

##### FLOW #####
Login -> Start New Operation -> Add sale* -> End Operation

##### ROUTES #####
1. Agent
```
Agent login | /ag/li/:user-:pass
Get list of all agents | /ag
Find agent by username | /ag/:user
```

2. Operation
```
Get all operations | /op
Start new operation | /op/start/:location-:agent_user/bal:start_bal_cash-:start_bal_qr/inv/:start_item_bal
End operation | /op/end/:id/bal/:end_bal_cash-:end_bal_qr/inv/:end_item_bal
```

3. Sales
```
Get all sales | /sa
Get list of sale (by date) | /sa/find/:syear-:smonth-:sday-:eyear-:emonth-:eday
Get list of sale (by operation_id) | /sa/find/:operation-id
New Sale | /sa/new/:amount-:qty-:payment_type-:operation_id-:item_id
Update Sale | 
Delete Sale (admin only / only the most recent one) | /sa/del/:id
```

4. Inventory
```
Get all inventory | /inv
Add inventory | /inv/new/:start_bal
Update inventory | /inv/up/:id-:end_bal
Delete inventory (admin only) |
```

##### NOTES #####
1. Payment types:-
    1 - Cash
    2 - QR Maybank
    3 - QR TnG
    99 - Free
2. Item id:-
    1 - MD2 Juice
    2 - MD2 Sliced Fruit
    3 - MD2 Raw Fruit
3. start_item_bal and end_item_bal data type is string so that I can format the data as `1=150&2=45`, 1 indicating item_id 1 and 2 indicating item_id 2
4. make sale_item_group table to group items when posting data (instead of having to make 3 post requests, for new sale)