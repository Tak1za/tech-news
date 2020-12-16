import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { makeStyles } from "@material-ui/core/styles";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import Typography from "@material-ui/core/Typography";
import Box from "@material-ui/core/Box";

function TabPanel(props) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`vertical-tabpanel-${index}`}
      aria-labelledby={`vertical-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box p={3}>
          <Typography component={"span"} noWrap={false}>
            {children}
          </Typography>
        </Box>
      )}
    </div>
  );
}

TabPanel.propTypes = {
  children: PropTypes.node,
  index: PropTypes.any.isRequired,
  value: PropTypes.any.isRequired,
};

function a11yProps(index) {
  return {
    id: `vertical-tab-${index}`,
    "aria-controls": `vertical-tabpanel-${index}`,
  };
}

const useStyles = makeStyles((theme) => ({
  root: {
    flexGrow: 1,
    backgroundColor: theme.palette.background.paper,
    display: "flex",
    justifyContent: "flex-start",
    height: "100%"
  },
  tabs: {
    borderRight: `1px solid ${theme.palette.divider}`,
  },
  listitem: {
    color: "black",
  },
}));

export default function VerticalTabs() {
  const classes = useStyles();
  const [value, setValue] = useState(0);
  const [hnData, setHnData] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8080/hn/stories")
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setHnData(data.data);
      });
  }, []);
  const [redditData, setRedditData] = useState([]);
  useEffect(() => {
    fetch("http://localhost:8080/r/stories")
      .then((res) => res.json())
      .then((data) => {
        console.log(data);
        setRedditData(data.data);
      });
  }, []);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  return (
    <div className={classes.root}>
      <Tabs
        orientation="vertical"
        variant="scrollable"
        value={value}
        onChange={handleChange}
        aria-label="Vertical tabs example"
        className={classes.tabs}
      >
        <Tab label="HackerNews" {...a11yProps(0)} />
        <Tab label="Reddit" {...a11yProps(1)} />
      </Tabs>
      <TabPanel value={value} index={0}>
        <List>
          {hnData.map((item) => {
            return (
              <ListItem
                key={item.id}
                component="a"
                href={item.url}
                target="blank"
                className={classes.listitem}
              >
                <Typography component={"span"}>
                  <ListItemText primary={item.title} secondary={item.url} />
                </Typography>
              </ListItem>
            );
          })}
        </List>
      </TabPanel>
      <TabPanel value={value} index={1}>
        <List>
          {redditData.map((item) => {
            return (
              <ListItem
                key={item.name}
                component="a"
                target="blank"
                href={item.url}
                className={classes.listitem}
              >
                <ListItemText primary={item.title} secondary={item.url} />
              </ListItem>
            );
          })}
        </List>
      </TabPanel>
    </div>
  );
}
