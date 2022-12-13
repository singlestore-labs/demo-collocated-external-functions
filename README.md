# Local external functions demo using Docker Compose

This repo includes a standalone demonstration of using [SingleStore][s2] [external functions][extfns]. External functions allow you to define a user defined function in SingleStore which, when executed, will call out to an externally defined web service.

[SingleStore][s2] is a scale out relational database optimized for transactions and analytics. This is just one of many extensibility options that SingleStore offers to help it run your workload. You can run SingleStore yourself or use our cloud offering ([sign up here for $500 in credits!][trial]).

## Using this repo

1. [Sign up][trial] for a SingleStore account and get your license key from [the customer portal][portal].
2. Add the license key to your shell environment:
   ```bash
   export SINGLESTORE_LICENSE="YOUR LICENSE KEY HERE"
   ```
3. Run `docker-compose up --build` in this directory, check the logs for errors
4. Open [SingleStore Studio in your browser][studio]
5. Login using the username `root` and password `test`
6. Navigate to the SQL Editor
7. Run the following SQL commands line by line:
   ```sql
   use test;

   -- let's first make sure it works
   select * from tokenize("Cool! External functions work!");

   -- we can extract just the tokens
   select encoded::tokens from tokenize("Cool! External functions work!");

   -- now let's join it with a table
   select id, encoded::tokens from sample_data, tokenize(sample_data.t);
   ```

[s2]: https://www.singlestore.com
[extfns]: https://docs.singlestore.com/db/latest/en/reference/sql-reference/procedural-sql-reference/create--or-replace--external-function.html
[trial]: https://www.singlestore.com/cloud-trial/
[portal]: https://portal.memsql.com/
[studio]: http://localhost:8080