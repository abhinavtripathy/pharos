const nlp = require('compromise')
nlp.extend(require('compromise-numbers'))
let math = require('mathjs')


let doc = nlp("there are $22 areas to improve and twenty thousand things")
let nums = doc.values().toNumber()
nums = nums.out('array')

console.log(math.evaluate('(2 * 3)'))