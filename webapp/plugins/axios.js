import https from 'https';

export default function ({
   $axios,
   store
                         }) {

  $axios.onResponseError((error) => {
    if (error.status === 401 && $auth.loggedIn && $auth.refreshToken.get()) {
      $auth.refreshTokens();
    }
  });

  const agent = new https.Agent({
    rejectUnauthorized: false
  });
  
  $axios.onRequest(config => {
    if (process.env.dev) {
      config.httpsAgent = agent;
    }
  });
}
