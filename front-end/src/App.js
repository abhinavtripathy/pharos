import React from 'react';
import { makeStyles, Grid, CssBaseline, Paper }  from '@material-ui/core';
import Draft from './Draft';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
  },

  grid: {
    padding: 10,
  },

  editor: {
    height: 500,
    padding: 20,
  },

  output: {
    height: 500,
    padding: 20,
  },

  terminal: {
    height: 150,
    padding: 20,
  },

  draft: {
    padding: 10,
  },
}));

function App() {
  const classes = useStyles();
  let out = "";

  function output(raw) {
    out = "";
    for (let i = 0; i < raw.blocks.length; ++i) {
      out += raw.blocks[i].text;
    }
    console.log(raw);
    console.log(out);
  }

  return (
    <div className={classes.root}>
      <CssBaseline/>
      <Grid container spacing={2} className={classes.grid}>
        <Grid item xs={12} md={8}>
          <Paper className={classes.editor} elevation={3}>
            <Draft output={output}/>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper className={classes.output} elevation={3}>
            {out}
          </Paper>
        </Grid>
        <Grid item xs={12} md={12}>
          <Paper className={classes.terminal} elevation={3}>
            Terminal
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

export default App;
