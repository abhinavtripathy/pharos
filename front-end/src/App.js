import React from 'react';
import { makeStyles, Grid, CssBaseline, Paper, AppBar, Toolbar, Typography, Button }  from '@material-ui/core';
import Draft from './Draft';
import Output from './Output';

const useStyles = makeStyles(theme => ({
  root: {
    flexGrow: 1,
    overflow: "hidden",
  },

  grid: {
    padding: 10,
  },

  editor: {
    height: "65vh",
    padding: 20,
  },

  output: {
    height: "65vh",
    padding: 20,
  },

  terminal: {
    height: "20vh",
    padding: 20,
  },

  draft: {
    padding: 10,
  },

  title: {
    flexGrow: 1,
  },
  button: {
    marginRight: theme.spacing(2),
  },
}));

function App() {
  const classes = useStyles();
  let out = "Output";

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
      <AppBar position="static">
            <Toolbar>
              <Typography variant="h6" className={classes.title}>
                Pharos
              </Typography>
              <Button color="secondary" variant="outlined" className={classes.button}>Help</Button>
              <Button color="secondary" variant="outlined">About</Button>
        </Toolbar>
      </AppBar>
      <Grid container spacing={2} className={classes.grid}>
        <Grid item xs={12} md={8}>
          <Paper className={classes.editor} elevation={3}>
            <Draft output={output}/>
          </Paper>
        </Grid>
        <Grid item xs={12} md={4}>
          <Paper className={classes.output} elevation={3}>
            <Output data={out}/>
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
