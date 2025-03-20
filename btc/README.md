# Test setup
Tests are using a local regression network. You need to start  
[merry](https://merry.dev/) to easily setup the testing environment. The `btctest`
package has some helper functions to use for the tests. 

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

