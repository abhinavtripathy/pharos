import React from 'react'

class Output extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            output: this.props.data,
        }
    }

    render() {
        console.log("props output: " + this.props.output);
        return (
            <div>
                {this.props.output}
            </div>
        )
    }
}

export default Output;