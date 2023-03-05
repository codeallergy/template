export default function ({
                           store,
                           redirect
                         }) {
  if (!store.state.auth.loggedIn) {
    return redirect('/auth/login')
  }
  const user = store.state.auth.user;
  if (user === undefined) {
    return redirect('/auth/login')
  }
  if (user.role !== 'ADMIN') {
    return redirect('/admin_required')
  }
}
