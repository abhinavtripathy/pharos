const nlp = require('compromise')
nlp.extend(require('compromise-numbers'))

let doc = nlp("")
  let nums = doc.values()
  console.log(nums.out('array'))