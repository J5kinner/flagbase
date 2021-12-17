import React from 'react';
import { render } from 'react-dom';
import { debugContextDevtool } from 'react-context-devtool';

import 'antd/dist/antd.css';

import Router from './router';
import Layout from 'antd/lib/layout/layout';
import { css, Global } from '@emotion/react';
import Poppins from '../fonts/Poppins-Regular.ttf';
import PoppinsBold from '../fonts/Poppins-Bold.ttf';
import PoppinsMedium from '../fonts/Poppins-Medium.ttf';

const mainElement = document.createElement('div');
mainElement.setAttribute('id', 'root');
document.body.appendChild(mainElement);
const container = document.getElementById('root');

const styles = css`
  @font-face {
    font-family: "Poppins";
    src: url(${Poppins}) format("truetype");
    font-weight: normal;
  font-style: normal;
  }

  @font-face {
    font-family: "Poppins";
    src: url(${PoppinsBold}) format("truetype");
    font-weight: bold;
  }

  @font-face {
    font-family: "Poppins";
    src: url(${PoppinsMedium}) format("truetype");
    font-weight: 500;
  }
`;

render(
  <Layout
    style={{ height: '100vh', padding: '50px', backgroundColor: 'white', fontFamily: 'Poppins' }}
  >

    <Global styles={styles} /> <Router />
  </Layout>,
  container
);

debugContextDevtool(container, {
  disable: process.env.NODE_ENV !== 'development'
});
