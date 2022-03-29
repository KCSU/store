import { useEditGroup } from "../../hooks/admin/useEditGroup";
import { useHasPermission } from "../../hooks/admin/useHasPermission";
import { Group } from "../../model/Group";
import { GroupDetailsForm } from "./GroupDetailsForm";
interface GroupProps {
  group: Group;
}

export function EditGroupForm({ group }: GroupProps) {
  const canWrite = useHasPermission("groups", "write");
  const mutation = useEditGroup(group.id);
  return (
    <GroupDetailsForm
      group={group}
      isDisabled={!canWrite}
      onSubmit={async (values) => {
        await mutation.mutateAsync(values);
      }}
    >
      Save Changes
    </GroupDetailsForm>
  );
}
