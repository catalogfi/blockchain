# Test setup
The test is mainly relying on local regression testnet. It's highly recommended to use 
[nigiri](https://nigiri.vulpem.com/) to easily set up a local testing environment. The `testutil`
package has some helper functions to use in the tests. You can check out the `btc.go` file in the 
`testutil` package for more details. 

There are some env variables need to be set before running all the tests. It's recommended to have 
a `.env` file and source it before running the tests. The test will first check if these variables 
and warn you if any of them is missing. Alternatively you can manually disable the envs check 
`btc_suite_test.go` if you just want to run a specific test which doesn't requiring any envs. 

# Test usage

```shell
ginkgo -v 

# Add the `cover` flag to get code coverage stats,
ginkgo -v --cover 
```

