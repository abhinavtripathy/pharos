
function interpreter(query) {

    let sentenceSplit = query.split(".")
    
    for(let i = 0; i < sentenceSplit.length; i++) {

        sentenceSplit[i] = sentenceSplit[i].trim()

    }
    
    sentenceSplit = sentenceSplit.filter(s => s != '')

    console.log("Sentence Split Array:")
    console.log(sentenceSplit)

    processQueries(sentenceSplit)

}

function processQueries(queryArr) {

    for(let query of queryArr) {
        
        console.log(query)

    }

}


function isMathQuery(query) {

    let mathQuery = query.split(" ")
    mathQuery = mathQuery.filter(s => s != '')

    let numCount = 0 
    let opCount = 0


    for(let mathQ of mathQuery) {

        if(!isNaN(parseInt(mathQ))) numCount++

        if(mathQ == '+' || mathQ == '-' || mathQ == '*' || mathQ == '/') opCount++

    }


    console.log(opCount)
    return numCount > 0 ? true : false && (opCount - 1) == numCount ? true : false

}

console.log(isMathQuery("3 +  4 / 2"))

test = "Hello world. This is a program and it is pretty good. I like it quite a bit a bit. a."

// interpreter(test)