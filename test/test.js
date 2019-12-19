let index = require('../index')

var assert = require('assert');
describe('Sum', function() {
  describe('Addition', function() {
    it('should return the sum', function() {
    
      assert.equal(index(1,2), 3);
    });
  });
});
