const nlp = require('compromise')
nlp.extend(require('compromise-numbers'))

let doc = nlp("there are 22 areas to improve and twenty thousand things")
let nums = doc.values().toNumber()
nums = nums.out('array')

  