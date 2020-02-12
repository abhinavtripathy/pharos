import React from 'react';
import {Editor, EditorState, convertToRaw} from 'draft-js';
import './editorStyles.css';
import { Button } from '@material-ui/core';

class Draft extends React.Component {
    constructor(props) {
        super(props);
        this.state = {editorState: EditorState.createEmpty()};;
        this.onChange = (editorState) => this.setState({editorState});
        this.setEditor = (editor) => {
          this.editor = editor;
        };
        this.focusEditor = () => {
          if (this.editor) {
            this.editor.focus();
          }
        };
        this.run = () => this.props.output(convertToRaw(this.state.editorState.getCurrentContent()));
    }
    
    componentDidMount() {
        this.focusEditor();
    }
    
    render() {
        return (
            <div id="root">
                <Button variant="contained" onClick={this.run}>
                    RUN
                </Button>
                <div id="editor" onClick={this.focusEditor}>
                    <Editor
                    ref={this.setEditor}
                    editorState={this.state.editorState}
                    onChange={this.onChange}
                    />
                </div>
            </div>
        );
    }
}

export default Draft;

