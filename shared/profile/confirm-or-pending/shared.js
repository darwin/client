// @flow
import type {Props} from '.'

const commonConfirm = ({platform, isPending}) => ({
  message: isPending
    ? 'Some proofs can take a few hours to recognize. Check back later.'
    : 'Leave your proof up so other users can identify you!',
  messageSubtitle: null,
  platformIcon: `icon-${platform}-logo-48`,
  platformIconOverlay: isPending ? 'icon-proof-pending' : 'icon-proof-success',
  title: isPending ? 'Your proof is pending.' : 'Verified!',
  usernameSubtitle: `@${platform}`,
})

export function propsForPlatform(props: Props): Object {
  switch (props.platform) {
    case 'twitter':
      return {
        ...commonConfirm(props),
      }
    case 'facebook':
      return {
        ...commonConfirm(props),
      }
    case 'reddit':
      return {
        ...commonConfirm(props),
      }
    case 'github':
      return {
        ...commonConfirm(props),
      }
    case 'hackernews':
      return {
        ...commonConfirm(props),
        message: props.isPending
          ? 'Hacker News caches its bios, so it might be a few hours before you can verify your proof. Check back later.'
          : 'Leave your proof up so other users can identify you!',
      }
    case 'dns':
      return {
        ...commonConfirm(props),
        message: props.isPending
          ? 'DNS proofs can take a few hours to recognize. Check back later.'
          : 'Leave your proof up so other users can identify you!',
      }
    case 'zcash':
      return {
        ...commonConfirm(props),
        message: 'Your Zcash address has now been signed onto your profile.',
        messageSubtitle: null,
        platformIcon: `icon-${props.platform}-logo-48`,
        platformIconOverlay: 'icon-proof-success',
        title: 'Verified!',
        usernameSubtitle: null,
      }
    case 'btc':
      return {
        ...commonConfirm(props),
        message: 'Your Bitcoin address has now been signed onto your profile.',
        messageSubtitle: null,
        platformIcon: `icon-${props.platform}-logo-48`,
        platformIconOverlay: 'icon-proof-success',
        title: 'Verified!',
        usernameSubtitle: null,
      }
    case 'http':
      return {
        ...commonConfirm(props),
        messageSubtitle: `Note: ${
          props.username
        } doesn't load over https. If you get a real SSL certificate (not self-signed) in the future, please replace this proof with a fresh one.`,
      }
    case 'https':
    case 'web':
      return {
        ...commonConfirm(props),
      }
    default:
      return {}
  }
}
