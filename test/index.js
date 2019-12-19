var sum = require('../index');
var expect = require('chai').expect;


it('should return sum of arguments', function() {
        expect(sum(1, 2)).to.equal(3)
  }
  

// describe('#sum()', function() {
  
//   context('with number arguments', function() {
//     it('should return sum of arguments', function() {
//       expect(sum(1, 2)).to.equal(3)
//     }) } 
    
//     )})


// var assert = require('assert');
// describe('Basic Mocha String Test', function () {
//  it('should return number of charachters in a string', function () {
//         assert.equal("Hello".length, 4);
//     });
//  it('should return first charachter of the string', function () {
//         assert.equal("Hello".charAt(0), 'H');
//     });
// });