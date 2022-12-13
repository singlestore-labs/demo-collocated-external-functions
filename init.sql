create database if not exists test;
use test;

create or replace external function tokenize (t TEXT)
returns TABLE (encoded TEXT)
as remote service "ext_fns:8000/text/tokenize"
format JSON;

create table sample_data (id int, t TEXT);
insert into sample_data values (1, 'This is a test');
insert into sample_data values (2, 'This is another test');
insert into sample_data values (3, 'This is a third test');