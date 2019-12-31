
function interpreter(query) {

    let sentenceSplit = query.split(".")
    
    for(let i = 0; i < sentenceSplit.length; i++) {

        sentenceSplit[i] = sentenceSplit[i].trim()

    }
    sentenceSplit = sentenceSplit.filter(v => v != '')

    console.log(sentenceSplit)

}

function processQueries(queryArr) {

    for(let query of queryArr) {

        

    }

}

test = "Hello world. This is a program and it is pretty good. I like it quite a bit a bit. a."

interpreter(test)