import { ComponentWithAs } from "@chakra-ui/react";
import { MotionProps, motion, isValidMotionProp } from "framer-motion";
import { forwardRef } from "react";

export type GenericMotionProps<Props> = Omit<Props, keyof MotionProps> &
  MotionProps & {
    as?: React.ElementType;
  };

type As<Props = any> = React.ElementType<Props>;

// HACK: allows us to use motion + transition + chakra
/**
 * This hack should be used instead of `motion<Props>(Component)` in cases where
 * we want to wrap Chakra components but also to use the
 * `transition` prop.
 * @example const MotionComponent = generateMotion<Props, 'div'>(Component);
 */
export function generateMotion<Props extends object, T extends As>(
  Component: React.ComponentClass<Props, any> | React.FunctionComponent<Props>
) {
  const Wrapped = motion<any>(
    forwardRef<Props, T>((props, ref) => {
      const chakraProps = Object.fromEntries(
        // do not pass framer props to DOM element
        Object.entries(props).filter(([key]) => !isValidMotionProp(key))
      );
      // @ts-ignore
      return <Component ref={ref} {...chakraProps} />;
    })
  ) as ComponentWithAs<T, GenericMotionProps<Props>>;
  Wrapped.displayName = `Motion${Component.displayName}`;
  return Wrapped;
}