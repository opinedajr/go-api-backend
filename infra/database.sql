create database gobank;
create table if not exists accounts (
	id UUID primary key,
	first_name varchar(50),
	last_name varchar(50),
	document varchar(50),
	agency serial,
	number serial,
	balance serial,
	created_at timestamp,
	modified_at timestamp NULL DEFAULT NULL
);

INSERT INTO accounts  (
id,
first_name,
last_name,
document,
agency,
number,
balance,
created_at,
modified_at
) 
VALUES (
	'4c082c46-80cd-482d-92cf-fae572bd2f65',
	'Jose',
	'Pineda',
	'00200200233',
	1,
	25236,
	0,
	'2023-10-22T15:04:05Z07:00',
	'2023-10-22T15:04:05Z07:00'
)

INSERT INTO accounts  (
id,
first_name,
last_name,
document,
agency,
number,
balance,
created_at,
modified_at
) 
VALUES (
	'8e823d44-27e5-4b34-a23f-0bb037f3c45b',
	'Jhon',
	'Snow',
	'12345678991',
	1,
	8907,
	0,
	'2023-10-22T15:04:05Z07:00',
	'2023-10-22T15:04:05Z07:00'
)