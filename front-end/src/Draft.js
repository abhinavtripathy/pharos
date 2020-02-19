import React from 'react';
import {Editor, EditorState, convertToRaw, CompositeDecorator} from 'draft-js';
import { Button } from '@material-ui/core';
import PlayArrowRoundedIcon from '@material-ui/icons/PlayArrowRounded';

import './editorStyles.css';

const vars = ["num", "bool", "string"];
const print_statements = ["out", "output", "print"];
const logic = ["if", "otherwise if", "otherwise", "then", "end if"];
const op_comp = ["is greater than", "is less than", "is equal to",];
const loop = ["loop", "until", "end loop"];
const functions = ["codeblock", "end codeblock"];

const VARS_REGEX = new RegExp("(?:[\\s]|^)(" + vars.join("|") + ")(?=[\\s]|$)", 'gi');
const PRINT_REGEX = new RegExp("(?:[\\s]|^)(" + print_statements.join("|") + ")(?=[\\s]|$)", 'gi')
const LOGIC_REGEX = new RegExp("(?:[\\s]|^)(" + logic.join("|") + ")(?=[\\s]|$)", 'gi')
const OPCOMP_REGEX = new RegExp("(?:[\\s]|^)(" + op_comp.join("|") + ")(?=[\\s]|$)", 'gi')
const LOOP_REGEX = new RegExp("(?:[\\s]|^)(" + loop.join("|") + ")(?=[\\s]|$)", 'gi')
const FUNC_REGEX = new RegExp("(?:[\\s]|^)(" + functions.join("|") + ")(?=[\\s]|$)", 'gi')

const DOUBLE_QUOTE_REGEX = /(["])(?:(?=(\\?))\2.)*?\1/g;
const SINGLE_QUOTE_REGEX = /(['])(?:(?=(\\?))\2.)*?\1/g;

const Variable = ({ children }) => {
    return (
        <span style={{ color:  "lightBlue"}}>
            {children}
        </span>
    );
};

const Print = ({ children }) => {
    return (
        <span style={{ color: "red"}}>
            {children}
        </span>
    );
};

const Logic = ({ children }) => {
  return (
      <span style={{ color: "mediumOrchid"}}>
          {children}
      </span>
  );
};

const OpComp = ({ children }) => {
  return (
      <span style={{ color: "mediumAquaMarine"}}>
          {children}
      </span>
  );
};

const Loop = ({ children }) => {
  return (
      <span style={{ color: "lightSalmon"}}>
          {children}
      </span>
  );
};

const Func = ({ children }) => {
  return (
      <span style={{ color: "gold"}}>
          {children}
      </span>
  );
};

const DoubleQuoteSpan = ({ children }) => {
    return (
      <span style={{ color: "orange" }}>
        {children}
      </span>
    );
};

const SingleQuoteSpan = ({ children }) => {
    return (
      <span style={{ color: "darkOrange" }}>
        {children}
      </span>
    );
};

function findWithRegex(regex, contentBlock, callback) {
    const text = contentBlock.getText();
    let matchArr, start;
    while ((matchArr = regex.exec(text)) !== null) {
      start = matchArr.index;
      callback(start, start + matchArr[0].length);
    }
}

function variableStrategy(contentBlock, callback) {
    findWithRegex(VARS_REGEX, contentBlock, callback);
}

function printStrategy(contentBlock, callback) {
    findWithRegex(PRINT_REGEX, contentBlock, callback);
}

function logicStrategy(contentBlock, callback) {
  findWithRegex(LOGIC_REGEX, contentBlock, callback);
}

function opcompStrategy(contentBlock, callback) {
  findWithRegex(OPCOMP_REGEX, contentBlock, callback);
}

function loopStrategy(contentBlock, callback) {
  findWithRegex(LOOP_REGEX, contentBlock, callback);
}

function funcStrategy(contentBlock, callback) {
  findWithRegex(FUNC_REGEX, contentBlock, callback);
}

function doubleQuoteStrategy(contentBlock, callback) {
    findWithRegex(DOUBLE_QUOTE_REGEX, contentBlock, callback);
}

function singleQuoteStrategy(contentBlock, callback) {
    findWithRegex(SINGLE_QUOTE_REGEX, contentBlock, callback);
}

const createDecorator = () =>
  new CompositeDecorator([
    {
        strategy: variableStrategy,
        component: Variable
    },
    {
        strategy: printStrategy,
        component: Print
    },
    {
      strategy: logicStrategy,
      component: Logic
    },
    {
      strategy: opcompStrategy,
      component: OpComp
    },
    {
      strategy: loopStrategy,
      component: Loop
    },
    {
      strategy: funcStrategy,
      component: Func
    },
    { 
        strategy: doubleQuoteStrategy,
        component: DoubleQuoteSpan
    },
    { 
        strategy: singleQuoteStrategy,
        component: SingleQuoteSpan
    },
]);

function Draft(props) {
  const [editorState, setEditorState] = React.useState(
    EditorState.createEmpty(createDecorator())
  );

  const editor = React.useRef(null);

  function focusEditor() {
    editor.current.focus();
  }

  function run() {
    props.output(convertToRaw(editorState.getCurrentContent()));
  }

  React.useEffect(() => {
    focusEditor()
  }, []);

  return (
    <div id="subroot">
        <Button variant="contained" color="secondary" onClick={run}>
            <PlayArrowRoundedIcon/>
        </Button>
        <div id="editor" onClick={focusEditor}>
            <Editor
            ref={editor}
            editorState={editorState}
            onChange={editorState => setEditorState(editorState)}
            />
        </div>
    </div>
  );
}

export default Draft;