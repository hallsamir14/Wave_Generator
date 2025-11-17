#TODO
(11/16/25):
- Update test suite for `generator` module to test for stateful generation
- Review `generator_test.go:134: case "triangle 4-sample period values (implementation-specific)"`: sample[0] = -32767, want 32767

(11/9/25):
- ~~Review and refine `generator.go` module~~

(11/5/25):
- ~~Implement stateful generator (retain state for each instance)~~
- ~~Refactor to seperate numeric samples and byte encoding~~
- ~~Implement `Generate` method for generator struct~~

(10/26/25)
- ~~Read generator test cases~~
- ~~Read wav_header test cases~~
- ~~Review and cleanup generator module~~
- ~~Cleanup main function~~

(10/19/25)
- ~~Continue implementing functionaly to generate different wave forms~~
- ~~Start test suite for generator module~~
- ~~Complete square wave implementation~~
- ~~Random suffix for each generated file - keep uniquely named samples~~ 

(10/18/25)
- ~~Review writer module test suites~~
- ~~Designate target path for written wav files~~
- ~~Review math for triangle wav form~~
- ~~Start implementing functionaly to generate different wave forms~~
